# Use an official Python
FROM resin/rpi-raspbian

# Set workdit to /app
WORKDIR /app

#copy current directoyr contents into app
ADD . /app

# install any needed packages specified in requirements
RUN apt-get update
RUN apt-get install gcc python python-dev python-pip
RUN pip install -r requirements.txt

# Make oirt 80 available to the outside world throught this container
EXPOSE 80

# Define enviroment variable
ENV NAME World

# Run app.py when container launches
CMD ["python","app.py"]
