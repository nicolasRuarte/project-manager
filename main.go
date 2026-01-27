package main
import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type project struct {
	Name string `json:"name"`
	Directory string `json:"directory"`
	ProjectType string `json:"projectType"` // If it is a web, mobile, dekstop or videogame project (Not limited to those options)
	CreationDate time.Time `json:"creationDate"`
	LastAccessed time.Time `json:"lastAccessed"`
}

const projectsFilePath = "projects.json"

func projectsFileExists() bool {
	_, err := os.ReadFile(projectsFilePath)
	if err != nil {
		return false
	}

	return true
}

func CreateProject() {
	fmt.Println("\nCreating project...")

	var name string
	var directory string
	var projectType string

	fmt.Print("Project name: ")
	fmt.Scan(&name)
	fmt.Print("Project directory: ")
	fmt.Scan(&directory)
	fmt.Print("Project type (web, dekstop, game, etc): ")
	fmt.Scan(&projectType)


	if !projectsFileExists() {
		// Using Unix time 0 as a way of implementing a null time. Fix later
		var projects []project
		newProject := project{name, directory, projectType, time.Now(), time.Unix(0, 0)}

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
	
	var savedProjects project
	err2 := json.Unmarshal(jsonData, &savedProjects)
	if err2 != nil {
		fmt.Println("Error: ", err2)
	}

	fmt.Println("Saved projects: ", savedProjects)
}

func readProject() {}

func updateProject() {}

func deleteProject() {}

func workInProject() {}

func main() {

	fmt.Println("Welcome to the project manager")
	
	fmt.Println("Please choose one of the options below: ")
	fmt.Println("	1) Create a project")
	fmt.Println("	2) Go to a project (WIP)")
	fmt.Println("	3) Update a project (WIP)")
	fmt.Println("	4) Delete a project (WIP)")
	//fmt.Println("	5) Work in project")
	fmt.Print("Your option: ")

	var input int 
	fmt.Scanln(&input)

	if input > 4 || input < 1 {
		fmt.Println("Error: The inserted value is not one of the permitted options")
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

	// Eliminar después de testear
	os.Remove("projects.json")
	fmt.Println("Borrando json")
}
