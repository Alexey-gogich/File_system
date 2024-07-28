package router

import (
	"Frank_rg/internal/controllers"
	"Frank_rg/internal/storage/database"

	"github.com/gorilla/mux"
)

func Routes(folders []*database.Folder) *mux.Router {
	mux := mux.NewRouter()
	controller := controllers.Controller{
		Folders: folders,
	}
	mux.HandleFunc("/conductor/{folder_name}", controller.Get_folder)

	//folders
	mux.HandleFunc("/conductor/{folder_name}/create_folder", controller.Create_folder)
	mux.HandleFunc("/conductor/{folder_name}/{child_name}/update_folder_name", controller.Update_folder_name)

	//files
	mux.HandleFunc("/conductor/{folder_name}/create_file", controller.Create_file)
	mux.HandleFunc("/conductor/{folder_name}/{file_name}/update_file_name", controller.Update_file_name)
	mux.HandleFunc("/conductor/{folder_name}/{file_name}/delete_file", controller.Delete_file)
	mux.HandleFunc("/conductor/{folder_name}/{file_name}/download_file", controller.Download_file)
	return mux
}
