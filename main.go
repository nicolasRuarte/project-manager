package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
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

func projectsFileExists() bool {
	_, err := os.ReadFile(projectsFilePath)
	if err != nil {
		return false
	}

	return true
}

func processYesOrNoInput(input string) bool {
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
	hasInitScript = processYesOrNoInput(yesOrNo)

	// Using Unix time 0 as a way of implementing a null time. Fix later
	newProject := project{name, directory, projectType, time.Now(), time.Unix(0, 0), hasInitScript}

	if !projectsFileExists() {
		var projects []project
		projects = append(projects, newProject)

		jsonBytes, err := json.Marshal(projects)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		err2 := os.WriteFile(projectsFilePath, jsonBytes, 0644)
		if err2 != nil {
			fmt.Println("Error: ", err2)
			return
		}

		//fmt.Println("Project created successfully!")
		
		_, err3 := os.ReadFile(projectsFilePath)
		if err3 != nil {
			fmt.Println("Error: ", err3)
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
	err2 := json.Unmarshal(jsonData, &savedProjects)
	if err2 != nil {
		fmt.Println("Error: ", err2)
	}

	savedProjects = append(savedProjects, newProject)
	
	jsonResult, err3 := json.Marshal(savedProjects)
	if err3 != nil {
		fmt.Println("Error: ", err3)
		return
	}
	err4 := os.WriteFile(projectsFilePath, jsonResult, 0644)
	if err4 != nil {
		fmt.Println("Error: ", err4)
		return
	}
}

func readProject() {
	jsonBytes, err := os.ReadFile(projectsFilePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var savedProjects []project
	err2 := json.Unmarshal(jsonBytes, &savedProjects)
	if err2 != nil {
		fmt.Println("Error: ", err2)
	}

	// Imprime los proyectos guardados en pantalla
	fmt.Println("\nSelect a project: ")
	for i := 0; i < len(savedProjects); i++ {
		if i == 0 {
			fmt.Printf("	%d) %s", i + 1, savedProjects[i].Name)
			continue
		}
		fmt.Printf("\n	%d) %s", i + 1, savedProjects[i].Name)
	}

	var selectedProjectId int

	fmt.Print("\nSelect the index of a project: ")
	_,  err3 := fmt.Scan(&selectedProjectId)
	if err3 != nil {
		log.Fatal(err3)
	}

	invalidInputRange := selectedProjectId < 0 || selectedProjectId > len(savedProjects)
	if invalidInputRange {
		fmt.Println("Inserted value is not allowed")
		return
	}

	// -1 porque los índices seleccionables empiezan de uno y la posición del array empieza desde el 0
	selectedProject := savedProjects[selectedProjectId - 1]
	fmt.Println("PROJECT: ")
	fmt.Println("Name: ", selectedProject.Name)
	fmt.Println("Directory: ", selectedProject.Directory)
	fmt.Println("Project type: ", selectedProject.ProjectType)
	fmt.Println("Creation date: ", selectedProject.CreationDate)
	fmt.Println("Last accessed: ", selectedProject.LastAccessed)
	fmt.Println("Has init script: ", selectedProject.HasInitScript)
}

func updateProject() {
	jsonData, err := os.ReadFile(projectsFilePath)
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}

	var savedProjects []project
	err2 := json.Unmarshal(jsonData, &savedProjects)
	if err2 != nil {
		log.Fatal("Error: ", err2)
		return
	}

	fmt.Println("\nSelect a project: ")
	for i := 0; i < len(savedProjects); i++ {
		if i == 0 {
			fmt.Printf("\t%d) %s", i + 1, savedProjects[i].Name)
			continue
		}
		fmt.Printf("\n\t%d) %s", i + 1, savedProjects[i].Name)
	}

	fmt.Print("\nInsert the index of the project: ")
	var selectedProjectIndex int
	fmt.Scan(&selectedProjectIndex)

	invalidInputRange := selectedProjectIndex < 1 || selectedProjectIndex > len(savedProjects)
	if invalidInputRange {
		fmt.Println("Inserted value is not allowed")
		return
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

	invalidInputRange = selectedAttributeIndex < 1 || selectedAttributeIndex > 4
	if invalidInputRange {
		fmt.Println("Inserted value is not allowed")
		return
	}

	switch selectedAttributeIndex {
	case 1:
		var newValue string
		fmt.Print("Insert new value: ")
		fmt.Scan(&newValue)

		savedProjects[selectedProjectIndex - 1].Name = newValue

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

		savedProjects[selectedProjectIndex - 1].Directory = newValue

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

		savedProjects[selectedProjectIndex - 1].ProjectType = newValue

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
		savedProjects[selectedProjectIndex - 1].HasInitScript = !savedProjects[selectedProjectIndex - 1].HasInitScript 

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

func deleteProject() {}

func workInProject() {}

func main() {

	fmt.Println("Welcome to the project manager")
	
	fmt.Println("Please choose one of the options below: ")
	fmt.Println("	1) Create a project")
	fmt.Println("	2) Go to a project")
	fmt.Println("	3) Update a project")
	fmt.Println("	4) Delete a project (WIP)")
	fmt.Println("	5) Work in project (WIP)")
	fmt.Print("Your option: ")

	var input int 
	fmt.Scanln(&input)

	if input > 4 || input < 1 {
		fmt.Println("Error: The inserted value is not one of the allowed options")
		return
	}

	switch input {
	case 1:
		CreateProject()
	case 2:
		readProject()
	case 3:
		updateProject()
	case 4: deleteProject()
	}
}
