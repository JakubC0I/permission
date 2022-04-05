package module

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	_        int = 1 << iota // 0 = no permissions
	Viewer                   //can view files
	Executer                 //can view and execute files
	Editor                   //can edit and view files
	Manager                  //can view, edit and execute files
	Finances                 //Finanzabteilung
	IT                       //ITabteilung
	Service                  //Service_Desk
)

type Variables struct {
	AddPfad      string
	RemovePfad   string
	AddUser      string
	Genehmigen   string
	AddImage     string
	Ticket       string
	Login        string
	Register     string
	Secret       string
	Notification string
	Deny         string
	UserID       string
}

type Success struct {
	Success bool               `bson:"success" json:"success"`
	Message string             `bson:"info" json:"info"`
	Info    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
}

type UserCookie struct {
	IDfromCookie string `json:"_id"`
	Data         struct{}
}

type UserArray struct {
	IDS   []primitive.ObjectID `bson:"_ids,omitempty" json:"_ids"`
	Pfads []string             `bson:"pfads" json:"pfads"`
}

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Email         string             `bson:"email,omitempty"`
	Password      string             `bson:"password"`
	Date          time.Time          `bson:"created at"`
	PermissionBit int32              `bson:"permissions" json:"permissions"`
	Pfads         []string           `bson:"pfads" json:"pfads"`
	Genehmiger    primitive.ObjectID `bson:"genehmiger" json:"genehmiger"`
}

type Task struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Besteller  primitive.ObjectID   `bson:"besteller" json:"besteller"`
	Action     string               `bson:"action" json:"action"`
	Genehmigt  bool                 `bson:"genehmigt" json:"genehmigt"`
	Genehmiger primitive.ObjectID   `bson:"genehmiger" json:"genehmiger"`
	Created_at time.Time            `bson:"created at"`
	Data       []string             `bson:"data" json:"data"`
	Betroffene []primitive.ObjectID `bson:"_bids" json:"_bids"`
}

//ODPOWIEDŹ dlaczego musi być bson i json przy int
//https://stackoverflow.com/questions/58075716/why-am-i-getting-all-zero-value-for-certain-field-in-my-json-from-mongodb
