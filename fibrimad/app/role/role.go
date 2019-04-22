package role

import (
	"github.com/nayimotis/raggaer/fibrimad/app/config"
	"github.com/nayimotis/raggaer/fibrimad/app/models"
)

const (
	// Login role that allows users to login into their accounts
	Login = "login"

	// CreateUser role that allows to create more users
	CreateUser = "create_user"

	// UserList role that allows to view the user list
	UserList = "user_list"

	// EditUser role that allows to edit users
	EditUser = "edit_user"

	// DeleteUser role that allows to delete users
	DeleteUser = "delete_user"

	// ViewUser role that allows to view users
	ViewUser = "view_user"

	// RemoveUser role that allows to remove users from work orders
	RemoveUser = "remove_user"

	// CreateWorkOrder role that allows to create a work order
	CreateWorkOrder = "create_work_order"

	// WorkOrderList role that allows to see the list of work orders
	WorkOrderList = "work_order_list"

	// DeleteWorkOrder role that allows to delete a work order
	DeleteWorkOrder = "delete_work_order"

	// EditWorkOrder role that allows to edit a work order
	EditWorkOrder = "edit_work_order"

	// ViewWorkOrder role that allows to see a work order
	ViewWorkOrder = "view_work_order"

	// UploadFileWorkOrder role that allows to upload files for a work order
	UploadFileWorkOrder = "upload_file_work_order"

	// ViewFile role that allows to view any file
	ViewFile = "view_file"

	// DeleteFile role that allows to delete any file
	DeleteFile = "delete_file"

	// EditFile role that allows to edit any file
	EditFile = "edit_file"

	// CreateBox role that allows to create a box
	CreateBox = "create_box"

	// ViewBox role that allows to view a box
	ViewBox = "view_box"

	// EditBox role that allows to edit a box
	EditBox = "edit_box"

	// DeleteBox role that allows to delete a box
	DeleteBox = "delete_box"

	// UploadFileBox role that allows to upload files for a box
	UploadFileBox = "upload_file_box"

	// ViewFileBox role that allows to view a box file
	ViewFileBox = "view_file_box"
)

// List defiens the application detailed role list
var List = map[string]DetailedRole{
	CreateBox: DetailedRole{
		Description: "Permite crear una nueva caja",
	},
	ViewFile: DetailedRole{
		Description: "Permite ver cualquier fichero",
	},
	DeleteFile: DetailedRole{
		Description: "Permite eliminar cualquier fichero",
	},
	Login: DetailedRole{
		Description: "Permite iniciar sesión",
	},
	CreateUser: DetailedRole{
		Description: "Permite crear usuarios",
	},
	UserList: DetailedRole{
		Description: "Permite ver la lista de usuarios",
	},
	EditUser: DetailedRole{
		Description: "Permite editar usuarios",
	},
	ViewUser: DetailedRole{
		Description: "Permite visualizar usuarios",
	},
	RemoveUser: DetailedRole{
		Description: "Permite retirar usuarios de una obra",
	},
	DeleteUser: DetailedRole{
		Description: "Permite eliminar usuarios",
	},
	CreateWorkOrder: DetailedRole{
		Description: "Permite crear nuevas obras",
	},
	WorkOrderList: DetailedRole{
		Description: "Permite ver la lista de obras",
	},
	DeleteWorkOrder: DetailedRole{
		Description: "Permite eliminar obras",
	},
	EditWorkOrder: DetailedRole{
		Description: "Permite editar obras",
	},
	ViewWorkOrder: DetailedRole{
		Description: "Permite visualizar obras",
	},
	UploadFileWorkOrder: DetailedRole{
		Description: "Permite subir archivos a una obra",
	},
	EditFile: DetailedRole{
		Description: "Permite editar archivos",
	},
	ViewBox: DetailedRole{
		Description: "Permite visualizar una caja",
	},
	DeleteBox: DetailedRole{
		Description: "Permite eliminar una caja",
	},
	EditBox: DetailedRole{
		Description: "Permite editar una caja",
	},
	UploadFileBox: DetailedRole{
		Description: "Permite subir fotos a una caja",
	},
	ViewFileBox: DetailedRole{
		Description: "Permite visualizar fotos de una caja",
	},
}

// DetailedRole defines a detailed role information
type DetailedRole struct {
	Description string
}

// UserHasRole checks if the given user has a role
func UserHasRole(cfg *config.Config, u *models.User, role string) bool {
	configRole, ok := cfg.Roles[u.Role]
	if !ok {
		return false
	}

	for _, r := range configRole {
		if r == role {
			return true
		}
	}

	return false
}
