# Project Manager

Project Manager is an interactive CLI application designed to manage all your local repositories. If you are the type of person that saves its coding projects on separate directorys that have no correlation to one another, then this is your solution.

NOTE: This project probably only works on Linux

## Features

* Allows you to create, read, update and delete a list of your projects
* Allows you to execute an initialization script to open all the apps that you need to work on the project

## How to use

1. Execute the command project-manager.
2. Then you'll have printed on screen a list of numerated options. Select the index of the option you want to use
3. Follow the process step by step of the option you selected

## How to execute a script

1. Before anything, you will need to have an executable file called init. That file is the one that is going to be executed by the program and it should have the instructions to open all the apps that you need to work on the project. Creating one is easy: here's an example:

```
firefox "chatgpt.com"
sleep 5

code ~/mis-apps/project-manager/
sleep 5

yt-music 
sleep 5

gnome-terminal -- bash -c 'cd  ~/mis-apps/project-manager/; nvim .;  bash'
```

This one opens Firefox and Visual Studio Code, but you can write something similar that follows your needs

2. Now you execute the program and select "Work on project"
3. Then, you select the project you want to work on and the app will execute the init script you made earlier

## Important

This program is on its first versions yet, so it has some limitations, specially with processing input. Those limitations are the following:

1. When creating or updating a project, your project name has to be one word only. If you need spaces you can use underscores or hyphens.
2. When creating or updating a project, its directory has to be an absolute path. That means that for now you can't use "~/"
