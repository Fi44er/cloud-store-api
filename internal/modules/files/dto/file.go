package file_dto

type CreateFolderRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"` // UUID родительской папки
}
