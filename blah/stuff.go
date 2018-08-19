//CCCP stands for Central Computation and Control Processor. It handles unifying the state of the
//Game between the two front ends and all of the bots on the field
package main

import (
"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
	"math"
)

const NUMBOTS = 2//Number of bots in the game
const PXWIDTH = 640
const PXHEIGHT = 480
const OUTPORT = "1917" // The port the bots will recieve commands from
const CAMPORT = "1918" //The port the camera will send its data down
const REPROGRAMON = false
const CAROLYNSERVER = "http://ec2-52-15-233-132.us-east-2.compute.amazonaws.com:3000/url/"
const MARGIN = 10

//Running main process
func main_full() {
	runtime.GOMAXPROCS(10)
	//Loads all the IP addresses of FEs, CCCP, and bots
	ips := getConfig("ips.txt")
	var lohost string
	lohost = ips.Cpu[0]
	log.Println(lohost)
//	uploaddomain(CAROLYNSERVER,palhost)
	//This is the starting port for the sheeps response
	inportstart := 2000
	var commandwg sync.WaitGroup

	// String to communicate out with bots
	outServerAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(lohost, OUTPORT))
	CheckError(err)

	//Initializing camera
	camip := "localhost"
	cam := initcamera(NUMBOTS, CAMPORT,camip)

	//Initializing idhash. used to get a certain sheep object from an id
	camToIdx := make(map[uint64]*Sheep)

	//Initializing sheep connections
	sheeps := make([]*Sheep, len(ips.Bot))
	for i, ip := range ips.Bot {
		sheeps[i] = initsheep(ip, lohost, uint16(inportstart+i))
		sheeps[i].commands.sheepF &= 0xEF
		sheeps[i].commands.servoAngle = 90
	}

	if REPROGRAMON {
		for _, sheep := range sheeps {
			sheep.commands.sheepF = SHEEPREPROGRAM
			sheep.sendCommands(outServerAddr)
			sheep.commands.sheepF &= 0xBF
		}
		wait := time.NewTimer(10*time.Second)
		<-wait.C
		panic("REPROGRAMMED")
	}
	//IDing process. How it works is that for each sheep, its light is turned on, a moment is waited
	//And then the position of all the bots is found. All the id's found are then iterated over, and if
	//the new id which would have popped up is then found. That sheep is marked as at that position,
	//And the sheep is then loaded into the camera id hash
	for _, sheep := range sheeps {
		found := false
		sheep.commands.sheepF |= SHEEPRST
		sheep.commands.sheepF |= SHEEPLIGHT
		sheep.sendCommands(outServerAddr)

		wait := time.NewTimer(time.Second*2)
		<-wait.C

		ids, xs, ys,orients := cam.getPos()
		for idx,id := range ids {
			_, inHash := camToIdx[id]
			if !inHash {
				sheep.currX = xs[idx]
				sheep.currY = ys[idx]
				sheep.pathhead.X = float64(sheep.currX)
				sheep.pathhead.Y = float64(sheep.currY)
				sheep.commands.camOrient = uint16(orients[idx])
				camToIdx[id] = sheep
				found = true
				break
			}
		}
		if !found {
			panic("SHEEP NOT FOUND IN INITIALIZATION PROCESS")
		}else{
			sheep.sendCommands(outServerAddr)
			sheep.commands.sheepF &= 0xFE
		}
		log.Println("IDing ids",ids)
	}

	//Initializing Frontend Server
	datawrite := MkChanDataWrite(ips.Fes, NUMBOTS,sheeps)
	//Setup server for the unity to make put requests too
	http.HandleFunc("/", http.HandlerFunc(datawrite.Botctrlserve))
	http.HandleFunc("/info", http.HandlerFunc(datawrite.APIServe))
	go http.ListenAndServe(numtoportstr(80), nil)

	//Wait for both Front ends to check in
	log.Println("Waiting For FE to Connect")
	for datawrite.gamestart != true {
		wait := time.NewTimer(time.Millisecond*10)
		<-wait.C
	}

	//Random constant used in game
	gamedone := false
	top := true

	log.Println("Entering Game")
	for gamedone == false {

		//Update the position of all of the bots
		ids, xs, ys,orients := cam.getPos()
		for i, id := range ids {
			sheep,found := camToIdx[id]
			if found {
				sheep.currX = xs[i]
				sheep.currY = ys[i]
				sheep.commands.camOrient = uint16(orients[i])
			}
		}
		//Using these updated positions, update the frontend interface to reflect that
		datawrite.FE.UpdateGndBots(sheeps, false, false)

		//using the frontend commands, prepare the command structure for each sheep
		for _, sheep := range sheeps{
			pat, _, fire, turretAngl,_ := datawrite.frInfo(sheep)
			//Fire or don't fire
			sheep.firing(fire,0.4)
			//Set servo angle
			sheep.setTurrAngle(turretAngl)
			if top {
				sheep.commands.sheepF |= SHEEPLIGHT
			}else{
				sheep.commands.sheepF &= 0xEF
			}
			//Get next point to travel too
			next := getNextPoint(sheep,pat,50)
			dist := euclidDist(next.X, float64(sheep.currX),next.Y, float64(sheep.currY))

			sheep.commands.relDesY = int16(next.Y - float64(sheep.currY))
			sheep.commands.relDesX = int16(next.X - float64(sheep.currX))
			des_angle :=math.Atan2(-float64(sheep.commands.relDesY),float64(sheep.commands.relDesX))*180/3.141
			if des_angle < 0{
				des_angle += 360
			}
			/*
			if sheep.currX < MARGIN {
				sheep.commands.relDesY = 0
				sheep.commands.relDesX = 40
			}
			if sheep.currX > PXWIDTH - MARGIN {
				sheep.commands.relDesY = 0
				sheep.commands.relDesX = -40
			}
			if sheep.currY < MARGIN {
				sheep.commands.relDesY = 40
				sheep.commands.relDesX = 0
			}
			if sheep.currY > PXHEIGHT - MARGIN{
				sheep.commands.relDesY = -40
				sheep.commands.relDesX = 0
			}
			*/
			log.Println("####################GAME STEP ##################################")
			log.Println("Servo Desired: ",turretAngl)
			log.Println("Servo Actual: ",sheep.commands.servoAngle)
			log.Println("Dist: ",dist)
			log.Println("Next: ",next)
			log.Println("SheepHealth: ",sheep.resp.health)
			log.Println("Sheep Pos:",sheep.currX,sheep.currY)
			log.Println("Sheep Orient:",sheep.commands.camOrient)
			log.Println("Trying To get to: ", sheep.commands.relDesX, sheep.commands.relDesY)
			log.Println("Path: ",pat)
			log.Println("PathHead: ",sheep.pathhead)
			log.Println("Fire?: ",fire)
		}


		//Send those commands down to the bots. this is done in a burst of threads
		commandwg.Add(NUMBOTS)
		for _, sheep := range sheeps {
			go sheep.recState(&commandwg)
			wait := time.NewTimer(time.Nanosecond * 500)
			<-wait.C
			sheep.sendCommands(outServerAddr)
		}
		//Wait until all of the bots respond
		wait := time.NewTimer(time.Millisecond*50)
		<-wait.C
		commandwg.Wait()

		if datawrite.gamestatus != 0 {
			gamedone = true
		}

	}

	log.Println("GAME COMPLETE")
}

//Code to test the frontend
func main_frontend() {

	//Initializing sheep connections
	sheeps := make([]*Sheep, NUMBOTS)
	ips := getConfig("ips.txt")
	for i := 0; i < NUMBOTS; i += 1 {
		sheeps[i] = initsheep("localhost", "localhost", uint16(1000))
	}
	sheeps[0].currX = 250
	sheeps[0].currY = 250
	sheeps[1].currX = 400
	sheeps[1].currY = 400
	sheeps[0].resp.health = 10
	sheeps[1].resp.health = 10
	//uploaddomain(CAROLYNSERVER,palhost)
	//Initializing Frontend Server
	datawrite := MkChanDataWrite(ips.Fes, NUMBOTS,sheeps)
	//Setup server for the unity to make put requests too
	http.HandleFunc("/", http.HandlerFunc(datawrite.Botctrlserve))
	http.HandleFunc("/info", http.HandlerFunc(datawrite.APIServe))
	go http.ListenAndServe(numtoportstr(80), nil)

	log.Println("Waiting For FE to Connect")
	for datawrite.gamestart != true {
		wait := time.NewTimer(time.Millisecond*10)
		<-wait.C
	}


	log.Println("Enter Game")
	for {
		datawrite.FE.UpdateGndBots(sheeps, false, false)

		_, _, fire, turretAngl,_ := datawrite.frInfo(sheeps[0])
		sheeps[0].firing(fire,1)
		sheeps[0].setTurrAngle(turretAngl)
		log.Println(sheeps[0].commands.servoAngle)
		_, _, fire, turretAngl,_ = datawrite.frInfo(sheeps[1])
		sheeps[1].firing(fire,1)
		sheeps[1].setTurrAngle(turretAngl)
		sheeps[1].currY = (sheeps[1].currY + 1)%500

		wait := time.NewTimer(time.Millisecond*50)
		<-wait.C
	}
}
//Code to test the camera
func main_camera() {
	cam := initcamera(5, "1918","localhost")
	for {
		log.Println(cam.getPos())
	}
}

func main(){
	// This works and strip "/static/" fragment from path
	main_full()
	//main_camera()
	//main_frontend()
}
