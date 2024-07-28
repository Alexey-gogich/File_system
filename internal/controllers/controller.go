package controllers

import (
	"Frank_rg/internal/storage/database"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

type Controller struct {
	Folders []*database.Folder
}

func current_folder(folders []*database.Folder, current_folder_name string) *database.Folder {
	//Find current position
	for _, folder := range folders {
		if folder.Name == current_folder_name {
			current_folder := folder
			return current_folder
		}
	}
	return nil
}

func (c *Controller) Get_folder(w http.ResponseWriter, r *http.Request) {
	folder_name := mux.Vars(r)["folder_name"]
	for _, Folder := range c.Folders {
		if Folder.Name == folder_name {
			temp, err := template.ParseFiles("internal/templates/Conductor.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			temp.Execute(w, Folder)
			return
		}
	}
	http.Error(w, "Данной папки не существует", http.StatusNotFound)
}

// folders
func (c *Controller) Create_folder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		current_folder_name := mux.Vars(r)["folder_name"]
		current_folder := current_folder(c.Folders, current_folder_name)
		if current_folder == nil {
			fmt.Println("Папка не найдена")
			return
		}

		//Create folder
		new_folder := database.NewFolder(r.FormValue("folder_name"), current_folder)
		current_folder.Children = append(current_folder.Children, new_folder)
		c.Folders = append(c.Folders, new_folder)
		http.Redirect(w, r, "/conductor/"+current_folder_name, http.StatusFound)
	} else {
		temp, err := template.ParseFiles("internal/templates/Create_folder.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, struct{ Folder_name string }{Folder_name: mux.Vars(r)["folder_name"]})
	}
}

func (c *Controller) Update_folder_name(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		current_folder_name := mux.Vars(r)["folder_name"]
		current_folder := current_folder(c.Folders, current_folder_name)
		if current_folder == nil {
			fmt.Println("Папка не найдена")
			return
		}

		//Update folder name
		for _, folder := range current_folder.Children {
			if folder.Name == mux.Vars(r)["child_name"] {
				folder.Name = r.FormValue("folder_name")
			}
		}
		http.Redirect(w, r, "/conductor/"+current_folder_name, http.StatusFound)
	} else {
		temp, err := template.ParseFiles("internal/templates/Update_name.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, struct {
			Folder_name string
			Child_name  string
		}{Folder_name: mux.Vars(r)["folder_name"], Child_name: mux.Vars(r)["child_name"]})
	}
}

// files
func (c *Controller) Create_file(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		current_folder_name := mux.Vars(r)["folder_name"]
		current_folder := current_folder(c.Folders, current_folder_name)
		if current_folder == nil {
			http.Error(w, "", http.StatusInternalServerError)
			fmt.Println("Папка не найдена")
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		download_file, err := os.Create("internal/storage/files_storage/" + handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer download_file.Close()

		if _, err := io.Copy(download_file, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		current_folder.Files = append(current_folder.Files, database.NewFile(handler.Filename))
		http.Redirect(w, r, "/conductor/"+current_folder_name, http.StatusFound)
	} else {
		temp, err := template.ParseFiles("internal/templates/Create_file.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, struct {
			Folder_name string
		}{Folder_name: mux.Vars(r)["folder_name"]})
	}
}

func (c *Controller) Update_file_name(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		current_folder_name := mux.Vars(r)["folder_name"]
		current_folder := current_folder(c.Folders, current_folder_name)
		if current_folder == nil {
			fmt.Println("Папка не найдена")
			return
		}

		//Update file name
		for _, file := range current_folder.Files {
			if file.Name == mux.Vars(r)["file_name"] {
				if err := os.Rename("internal/storage/files_storage/"+file.Name, "internal/storage/files_storage/"+r.FormValue("file_name")); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				file.Name = r.FormValue("file_name")
				break
			}
		}
		http.Redirect(w, r, "/conductor/"+current_folder_name, http.StatusFound)
	} else {
		temp, err := template.ParseFiles("internal/templates/Update_name.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, struct {
			Folder_name string
			File_name   string
			Child_name  string
		}{Folder_name: mux.Vars(r)["folder_name"], File_name: mux.Vars(r)["file_name"]})
	}

}

func (c *Controller) Download_file(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		current_folder_name := mux.Vars(r)["folder_name"]
		current_folder := current_folder(c.Folders, current_folder_name)
		if current_folder == nil {
			fmt.Println("Папка не найдена")
			return
		}

		for _, file := range current_folder.Files {
			if file.Name == mux.Vars(r)["file_name"] {
				w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				http.ServeFile(w, r, "internal/storage/files_storage/"+file.Name)
			}
		}
		http.Error(w, "", http.StatusNotFound)
	}
}

func (c *Controller) Delete_file(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		current_folder_name := mux.Vars(r)["folder_name"]
		current_folder := current_folder(c.Folders, current_folder_name)
		if current_folder == nil {
			fmt.Println("Папка не найдена")
			return
		}

		//Delete file
		for index, file := range current_folder.Files {
			if file.Name == mux.Vars(r)["file_name"] {
				err := os.Remove("internal/storage/files_storage/" + mux.Vars(r)["file_name"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				current_folder.Files = append(current_folder.Files[:index], current_folder.Files[index+1:]...)
				http.Redirect(w, r, "/conductor/"+current_folder_name, http.StatusFound)
				break
			}
		}
	}
}
