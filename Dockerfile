# Use an official Python
FROM python:2.7-slim

# Set workdit to /app
WORKDIR /app

#copy current directoyr contents into app
ADD . /app

# install any needed packages specified in requirements
RUN pip install --trusted-host pypi.python.org -r requirements.txt

# Make oirt 80 available to the outside world throught this container
EXPOSE 80

# Define enviroment variable
ENV NAME World

# Run app.py when container launches
CMD ["python","app.py"]
