package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"mvdan.cc/sh/shell"
)

// Los campos a los que puede acceder JSON son los que tienen mayúscula nomás
type project struct {
	Name string `json:"name"`
	Directory string `json:"directory"`
	ProjectType string `json:"projectType"` // If it is a web, mobile, dekstop or videogame project (Not limited to those options)
	CreationDate time.Time `json:"creationDate"`
	LastAccessed time.Time `json:"lastAccessed"`
	HasInitScript bool `json:"hasInitScript"`
}

const projectsFilePath = "projects.json"

func CheckIfProjectsFileExists() bool {
	_, err := os.ReadFile(projectsFilePath)
	if err != nil {
		return false
	}

	return true
}

// Acá hay error porque me devuelve falso en casos que me debería devolver error. Corregir
func ProcessYesOrNoInput(input string) bool {
	if input != "y" && input != "Y" &&  input != "n" && input != "N" {
		fmt.Println("Please select 'y' or 'n' as options")
		return false
	}

	inputIsYes := input == "y" || input == "Y"
	if inputIsYes {
		return true
	} else {
		return false
	}
}

func GetProjectListFromJson() ([]project, error) {
	jsonBytes, err := os.ReadFile(projectsFilePath)
	if err != nil {
		return nil, err
	}
	
	var projectList []project
	err = json.Unmarshal(jsonBytes, &projectList)
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

// Devuelve el índice del proyecto elegido
func ShowSelectProjectMenu(savedProjects []project) (int, error) {
	const errorIntValue = -1
	fmt.Println("\nSelect a project: ")
	for i, project := range savedProjects {
		if i == 0 {
			fmt.Printf("	%d) %s", i + 1, project.Name)
			continue
		}
		fmt.Printf("\n	%d) %s", i + 1, project.Name)
	}

	var selectedProjectId int
	fmt.Print("\nSelect the index of a project: ")
	count,  err := fmt.Scan(&selectedProjectId)
	if err != nil {
		return errorIntValue, err
	}

	invalidInputRange := selectedProjectId < 1 || selectedProjectId > len(savedProjects) || count != 1
	if invalidInputRange {
		return errorIntValue, errors.New("El índice que ingresó es inválido")
	}

	// Esto porque en la interfaz la lista empieza en 1, mientras que los arrays en Go empiezan en 0
	projectIndexOnArray := selectedProjectId - 1
	
	return projectIndexOnArray, nil
}

func CreateProject() {
	fmt.Println("\nCreating project...")

	var name string
	var directory string
	var projectType string
	var hasInitScript bool
	var yesOrNo string // Variable para aceptar el y/n del usuario y transformarlo en booleano

	// Fix: El scan no acepta inputs que contengan espacios, hay que ver qué hacemos con eso
	fmt.Print("Project name: ")
	wordCount, err := fmt.Scan(&name)
	if wordCount > 1 {
		fmt.Println("Projects names should not have spaces")
	}
	fmt.Print("Project directory: ")
	fmt.Scan(&directory) // TODO: Agregar verificación de que exista el directorio
	fmt.Print("Project type (web, dekstop, game, etc): ")
	fmt.Scan(&projectType)
	fmt.Print("Has initialization script (y/n): ")
	fmt.Scan(&yesOrNo)
	hasInitScript = ProcessYesOrNoInput(yesOrNo)

	// Using Unix time 0 as a way of implementing a null time. Fix later
	newProject := project{name, directory, projectType, time.Now(), time.Unix(0, 0), hasInitScript}

	if !CheckIfProjectsFileExists() {
		var projects []project
		projects = append(projects, newProject)

		jsonBytes, err := json.Marshal(projects)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		err = os.WriteFile(projectsFilePath, jsonBytes, 0644)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		//fmt.Println("Project created successfully!")
		
		_, err = os.ReadFile(projectsFilePath)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		return
	}

	jsonData, err := os.ReadFile(projectsFilePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	
	var savedProjects []project
	err = json.Unmarshal(jsonData, &savedProjects)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	savedProjects = append(savedProjects, newProject)
	
	jsonResult, err := json.Marshal(savedProjects)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = os.WriteFile(projectsFilePath, jsonResult, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func ReadProject() {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)

	// -1 porque los índices seleccionables empiezan de uno y la posición del array empieza desde el 0
	selectedProject := savedProjects[projectIndex]
	fmt.Println("PROJECT: ")
	fmt.Println("Name: ", selectedProject.Name)
	fmt.Println("Directory: ", selectedProject.Directory)
	fmt.Println("Project type: ", selectedProject.ProjectType)
	fmt.Println("Creation date: ", selectedProject.CreationDate)
	fmt.Println("Last accessed: ", selectedProject.LastAccessed)
	fmt.Println("Has init script: ", selectedProject.HasInitScript)
}

func UpdateProject() {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Refactorizar, creo que puedo hacerlo más programático
	fmt.Println("\nSelect the attribute you want to update")
	fmt.Println("\t1) Name")
	fmt.Println("\t2) Directory")
	fmt.Println("\t3) Project type")
	fmt.Println("\t4) Has init script")

	var selectedAttributeIndex int
	fmt.Print("\nIndex of the attribute you want to update: ")
	fmt.Scan(&selectedAttributeIndex)

	invalidInputRange := selectedAttributeIndex < 1 || selectedAttributeIndex > 4
	if invalidInputRange {
		fmt.Println("Inserted value is not allowed")
		return
	}

	switch selectedAttributeIndex {
	case 1:
		var newValue string
		fmt.Print("Insert new value: ")
		fmt.Scan(&newValue)

		savedProjects[projectIndex].Name = newValue

		jsonBytes, err := json.Marshal(savedProjects)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		err = os.WriteFile(projectsFilePath, jsonBytes, 0644)
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 2:
		var newValue string

		fmt.Print("Insert new value: ")
		fmt.Scan(&newValue)

		savedProjects[projectIndex].Directory = newValue

		jsonBytes, err := json.Marshal(savedProjects)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		err = os.WriteFile(projectsFilePath, jsonBytes, 0644)
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 3:
		var newValue string

		fmt.Print("Insert new value: ")
		fmt.Scan(&newValue)

		savedProjects[projectIndex].ProjectType = newValue

		jsonBytes, err := json.Marshal(savedProjects)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		err = os.WriteFile(projectsFilePath, jsonBytes, 0644)
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 4: 
		savedProjects[projectIndex].HasInitScript = !savedProjects[projectIndex].HasInitScript 

		jsonBytes, err := json.Marshal(savedProjects)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		err = os.WriteFile(projectsFilePath, jsonBytes, 0644)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		fmt.Println("Has init script value was swapped")
	}

	fmt.Println("Project updated successfully!")
}

func DeleteProject() {
	fmt.Println("You can't delete a project yet. Functionality wasn't coded")
}

func WorkInProject() {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Hacer una verifcación de que existe el archivo init en el directorio del proyecto
	projectToWorkOn := savedProjects[projectIndex]

	fileInfo, err := os.Stat(projectToWorkOn.Directory + "init")
	if err != nil {
		log.Fatal("Error: ", err)
		fmt.Println("Check if you have an init file in your project's directory")
	}

	initFileIsExecutable := fileInfo.Mode()&0100 != 0

	if !initFileIsExecutable {
		log.Fatal("The init file in the project's directory is not executable")
	}

	fmt.Println("Executing script...")
	_, err = shell.SourceFile(context.TODO(), projectToWorkOn.Directory + "init")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println("Script executed. Good luck!")
}

func main() {

	fmt.Println("Welcome to the project manager")
	
	fmt.Println("Please choose one of the options below: ")
	fmt.Println("	1) Create a project")
	fmt.Println("	2) Go to a project")
	fmt.Println("	3) Update a project")
	fmt.Println("	4) Delete a project")
	fmt.Println("	5) Work in project")
	fmt.Print("Your option: ")

	var input int 
	fmt.Scanln(&input)

	if input > 5 || input < 1 {
		log.Fatal("Error: The inserted value is not one of the allowed options")
		return
	}

	switch input {
	case 1:
		CreateProject()
	case 2:
		ReadProject()
	case 3:
		UpdateProject()
	case 4: 
		DeleteProject()
	case 5:
		WorkInProject()
	}
}
