import json
with open("stringout.json","r+") as fh:
        jsoned = json.loads(fh.read())
        print json.dumps(jsoned,indent=4,sort_keys=True)
print jsoned["pull_request"]["head"]["ref"]
