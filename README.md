# Eskom Loadshed Forecasts

Loadshed is a microservice that can be used to get the current loadshedding forecasts for ESKOM.

## Installation

Download the latest release file for your operating system from https://github.com/Brumawen/loadsched/releases 

Extract the files and subfolders to a folder and run the following from the command line

        loadsched -service install
        loadsched -service run

This will install and run the loadsched microservice as a background service on your machine.

## Configuration

Once the microservice is running, navigate to http://localhost:20515/config.html in a web browser.

### EskomSePush Provider

You need to [subscribe to the EskomSePush API](https://eskomsepush.gumroad.com/l/api).

Once you have subscribed, an email will be sent to you with the Authentication Token.

Use the [Area Search (Text) API call](https://documenter.getpostman.com/view/1296288/UzQuNk3E#1986b098-ad88-436c-a5cd-5aa406e2fcf2)
to find the Area ID for your home.

Enter both the Area ID and the Token into the configuration fields and click Save.
