package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"mvdan.cc/sh/shell"
)

// Los campos a los que puede acceder JSON son los que tienen mayúscula nomás
type project struct {
	Name string `json:"name"`
	Directory string `json:"directory"`
	Tags []string `json:"tags"`
	CreationDate time.Time `json:"creationDate"`
	LastAccessed time.Time `json:"lastAccessed"`
	HasInitScript bool `json:"hasInitScript"`
}

// Path to the JSON file where projects are stored
var projectsFilePath string

/*
Checks if project-manager directory and projects.json file exists. 

If directory and/or file do not exist, it creates them
*/
func CheckIfProjectsFileExists() (bool, error) {
	homeDir, err := os.UserHomeDir()
	appDirectory := homeDir + "/.project-manager/"

	_, err = os.ReadDir(appDirectory)
	if err != nil {
		err = os.Mkdir(appDirectory, 0700)
		if err != nil {
			return false, errors.New("Error al crear la carpeta .project-manager/")
		}
	}

	_, err = os.ReadFile(projectsFilePath)
	if err != nil {
		err = os.WriteFile(projectsFilePath, nil, 0644)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func GetProjectListFromJson() ([]project, error) {
	jsonBytes, err := os.ReadFile(projectsFilePath)
	if err != nil {
		return nil, errors.New("File projects.json does not exist. Try creating a new project instead")
	}
	
	var projectList []project
	err = json.Unmarshal(jsonBytes, &projectList)
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func WriteToJsonFile(projects []project) error {
		jsonBytes, err := json.Marshal(projects)
		if err != nil {
			return err
		}

		// No sé bien qué hace esta variable, pero me deja usar la función
		const permissions = 0644
		err = os.WriteFile(projectsFilePath, jsonBytes, permissions)
		if err != nil {
			return err
		}

		return nil
}

// Procesa los strings 'y' y 'n' como true y false, respectivamente
func ProcessYesOrNoInput(input string) (bool, error) {
	if input != "y" && input != "Y" &&  input != "n" && input != "N" {
		return false, errors.New("Please select 'y' or 'n' as options") 
	}

	inputIsYes := input == "y" || input == "Y"
	if inputIsYes {
		return true, nil
	} else {
		return false, nil
	}
}


// Devuelve el índice del array en el que el proyecto está almacenado. Básicamente devuelve: índice de la opción de la interfaz - 1
func ShowSelectProjectMenu(savedProjects []project) (int, error) {
	const errorIntValue = -1
	fmt.Println("\nSelect a project: ")

	isArrayEmpty := len(savedProjects) == 0
	if isArrayEmpty {
		return errorIntValue, errors.New("No existe ningún proyecto. Intenta crear uno antes")
	}

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

func CreateProject() error {
	fmt.Println("\nCreating project...")

	fmt.Print("Project name: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)

	fmt.Print("Project directory: ")
	directory, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	directory = strings.TrimSpace(directory)
	const slashAsciiCode = 47
	lastDirCharacterIsNotSlash := directory[len(directory) - 1] != slashAsciiCode
	if lastDirCharacterIsNotSlash {
		var sb strings.Builder
		sb.WriteString(directory)
		sb.WriteString("/")
		directory = sb.String()
	}
	_, err = os.ReadDir(directory) // ReadDir() falla si no existe el directorio
	if err != nil {
		return err
	}

	fmt.Print("Tags: ")
	tagsInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	tagsInput = strings.TrimSpace(tagsInput)
	tags := strings.Split(tagsInput, " ")

	fmt.Print("Has initialization script (y/n): ")
	yesOrNo, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	yesOrNo = strings.TrimSpace(yesOrNo)
	hasInitScript, err := ProcessYesOrNoInput(yesOrNo)
	if err != nil {
		return err
	}

	// Using Unix time 0 as a way of implementing a null time. Fix later
	newProject := project{name, directory, tags, time.Now(), time.Unix(0, 0), hasInitScript}

	projectFileExists, err := CheckIfProjectsFileExists()
	if err != nil {
		return err
	}

	if !projectFileExists {
		var projects []project
		projects = append(projects, newProject)

		err = WriteToJsonFile(projects)
		if err != nil {
			return err
		}

		fmt.Println("Project created successfully!")

		return nil
	}

	jsonData, err := os.ReadFile(projectsFilePath)
	if err != nil {
		return err
	}
	
	var savedProjects []project
	err = json.Unmarshal(jsonData, &savedProjects)
	if err != nil {
		return err
	}

	savedProjects = append(savedProjects, newProject)
	
	err = WriteToJsonFile(savedProjects)
	if err != nil {
		return err
	}

	fmt.Println("Project created successfully!")

	return nil
}

func ReadProject() error {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		return err
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)
	if err != nil {
		return err
	}

	// -1 porque los índices seleccionables empiezan de uno y la posición del array empieza desde el 0
	selectedProject := savedProjects[projectIndex]
	fmt.Println("PROJECT: ")
	fmt.Println("Name: ", selectedProject.Name)
	fmt.Println("Directory: ", selectedProject.Directory)
	fmt.Println("Tags: ", selectedProject.Tags)
	fmt.Println("Creation date: ", selectedProject.CreationDate)
	if selectedProject.LastAccessed.Equal(time.Unix(0, 0)) {
		fmt.Println("Last accessed: Never")
	} else {
		fmt.Println("Last accessed: ", selectedProject.LastAccessed)
	}
	fmt.Println("Has init script: ", selectedProject.HasInitScript)

	return nil
}

func UpdateProject() error {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		return err
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)
	if err != nil {
		return err
	}

	// Refactorizar, creo que puedo hacerlo más programático
	fmt.Println("\nSelect the attribute you want to update")
	fmt.Println("\t1) Name")
	fmt.Println("\t2) Directory")
	fmt.Println("\t3) Tags")
	fmt.Println("\t4) Has init script")

	var selectedAttributeIndex int
	fmt.Print("\nIndex of the attribute you want to update: ")
	fmt.Scan(&selectedAttributeIndex)

	invalidInputRange := selectedAttributeIndex < 1 || selectedAttributeIndex > 4
	if invalidInputRange {
		return errors.New("Inserted value is not allowed")
	}

	reader := bufio.NewReader(os.Stdin)

	switch selectedAttributeIndex {
	case 1:
		fmt.Print("Insert new value: ")
		newValue, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		newValue = strings.TrimSpace(newValue)

		savedProjects[projectIndex].Name = newValue
	case 2:
		var newValue string

		fmt.Print("Insert new value: ")
		fmt.Scan(&newValue)

		savedProjects[projectIndex].Directory = newValue
	case 3:
		fmt.Print("Insert new value: ")
		newValue, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		newValue = strings.TrimSpace(newValue)
		tags := strings.Split(newValue, " ")

		savedProjects[projectIndex].Tags = tags
	case 4: 
		savedProjects[projectIndex].HasInitScript = !savedProjects[projectIndex].HasInitScript 

		fmt.Println("Has init script value was swapped")
	}

	err = WriteToJsonFile(savedProjects)
	if err != nil {
		return err
	}

	fmt.Println("Project updated successfully!")
	return nil
}

func DeleteProject() error {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		return err
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)
	if err != nil {
		return err
	}

	indexIsInvalid := projectIndex < 0 || projectIndex > len(savedProjects) - 1
	if indexIsInvalid {
		return errors.New("The selected option is invalid")
	}

	selectedIndexIsLastElement := projectIndex == len(savedProjects) - 1

	if selectedIndexIsLastElement {
		newProjectsArray := savedProjects[:projectIndex]

		err = WriteToJsonFile(newProjectsArray)
		if err != nil {
			return err
		}

		fmt.Println("Project deleted successfully!")

		return nil
	}

	newProjectsArray := append(savedProjects[:projectIndex], savedProjects[projectIndex + 1:]...)

	fmt.Println(newProjectsArray)
	err = WriteToJsonFile(newProjectsArray)
	if err != nil {
		return err
	}

	fmt.Println("Project deleted successfully!")

	return nil
}

func WorkInProject() error {
	savedProjects, err := GetProjectListFromJson()
	if err != nil {
		return err
	}

	projectIndex, err := ShowSelectProjectMenu(savedProjects)
	if err != nil {
		return err
	}

	// Hacer una verifcación de que existe el archivo init en el directorio del proyecto
	projectToWorkOn := savedProjects[projectIndex]

	fileInfo, err := os.Stat(projectToWorkOn.Directory + "init")
	if err != nil {
		fmt.Println("Check if you have an init file in your project's directory")
		return err
	}

	initFileIsExecutable := fileInfo.Mode()&0100 != 0
	if !initFileIsExecutable {
		return errors.New("The init file in the project's directory is not executable")
	}

	savedProjects[projectIndex].LastAccessed = time.Now()
	err = WriteToJsonFile(savedProjects)
	if err != nil {
		return err
	}

	fmt.Println("Executing script...")
	_, err = shell.SourceFile(context.TODO(), projectToWorkOn.Directory + "init")
	if err != nil {
		return err
	}
	fmt.Println("Script executed. Good luck!")

	return nil
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	projectsFilePath = homeDir + "/.project-manager/v0.2-projects.json"

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

	switch input {
	case 1:
		err = CreateProject()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 2:
		err = ReadProject()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 3:
		err = UpdateProject()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 4: 
		err = DeleteProject()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	case 5:
		err = WorkInProject()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	default:
		log.Fatal("Error: Inserted value is not allowed")
	}
}
