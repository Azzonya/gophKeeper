// Package tui provides a text-based user interface (TUI) for the GophKeeper client,
// allowing users to register, log in, and manage their data items through a simple
// interactive interface using tview.
package tui

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gophKeeper/client/internal/client"
	proto "gophKeeper/pkg/proto/gophkeeper"
)

// TUI represents the text-based user interface for the GophKeeper client,
// handling the display and interaction logic for user registration, login,
// and data item management.
type TUI struct {
	client *client.GophKeeperClient
	cache  *redis.Client
	app    *tview.Application
}

// NewTUI creates a new TUI instance with the given gRPC client, initializing
// the application and setting up the user interface.
func NewTUI(client *client.GophKeeperClient, rDB *redis.Client) *TUI {
	return &TUI{
		client: client,
		cache:  rDB,
		app:    tview.NewApplication(),
	}
}

func generateUniqueID() string {
	return uuid.New().String()
}

// Run starts the TUI application, displaying the main form with options for
// user registration, login, and quitting the application.
func (t *TUI) Run() error {
	form := tview.NewForm()

	form.AddButton("Register", t.register).
		AddButton("Login", t.login).
		AddButton("Quit", func() {
			t.app.Stop()
		})

	t.app.SetRoot(form, true)
	return t.app.Run()
}

// register handles the user registration process, displaying a form to input
// a username and password and sending the registration request to the server.
func (t *TUI) register() {
	if !t.client.ServerAvailable {
		return
	}
	form := tview.NewForm()

	form.
		AddInputField("Username", "", 20, nil, nil).
		AddInputField("Password", "", 20, nil, nil).
		AddButton("Register", func() {
			username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			_, err := t.client.Register(ctx, &proto.RegisterRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				t.showMessage("Register failed. Press Enter to go back.", t.restart)
				return
			}

			t.showMessage("Register successful. Press Enter to go back.", t.restart)
		}).
		AddButton("Cancel", func() {
			t.restart()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// login handles the user login process, displaying a form to input a username
// and password and sending the login request to the server.
func (t *TUI) login() {
	form := tview.NewForm()
	form.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs())

			resp, err := t.client.Login(ctx, &proto.LoginRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				t.showMessage("Login failed. Press Enter to go back.", t.restart)
				return
			}

			t.client.BearerToken = resp.Token

			t.showMessage("Login successful. Press Enter to open Menu.", t.showMainMenu)
		}).
		AddButton("Cancel", func() {
			t.restart()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// showMainMenu displays the main menu with options for creating, getting,
// updating, and deleting data items, as well as quitting the application.
func (t *TUI) showMainMenu() {
	menu := tview.NewList().
		AddItem("Create Data", "Create new data", 'c', t.createData).
		AddItem("Get Data", "Get existing data", 'g', t.getData).
		AddItem("List Data", "List existing data", 'l', t.listData).
		AddItem("Update Data", "Update existing data", 'u', t.updateData).
		AddItem("Delete Data", "Delete existing data", 'd', t.deleteData).
		AddItem("Quit", "Press to exit", 'q', func() {
			t.app.Stop()
		})

	t.app.SetRoot(menu, true).SetFocus(menu)
}

// createData displays a form for creating a new data item, allowing the user
// to input the type, data, and metadata, and sending the create request to the server.
func (t *TUI) createData() {
	if !t.client.ServerAvailable {
		t.showMessage("Server not available. Press Enter to go back.", t.showMainMenu)
		return
	}

	form := tview.NewForm()

	form.
		AddDropDown("Type", []string{"binary", "text", "credentials", "bank card"}, 0, nil).
		AddInputField("Data", "", 20, nil, nil).
		AddInputField("Meta", "", 20, nil, nil).
		AddButton("Submit", func() {
			_, typeField := form.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()
			dataField := form.GetFormItemByLabel("Data").(*tview.InputField).GetText()
			metaField := form.GetFormItemByLabel("Meta").(*tview.InputField).GetText()

			var data []byte
			var err error

			if typeField == "binary" {
				data, err = os.ReadFile(dataField)
				if err != nil {
					t.showMessage(fmt.Sprintf("Failed to read file: %v", err), t.showMainMenu)
					return
				}
			} else {
				data = []byte(dataField)
			}

			req := &proto.CreateDataRequest{
				Data: &proto.DataItem{
					Id:   generateUniqueID(),
					Type: typeField,
					Data: data,
					Meta: metaField,
				},
			}

			ctx, cancel := t.client.CreateContextWithMetadata(15 * time.Second)
			defer cancel()

			_, err = t.client.CreateData(ctx, req)
			if err != nil {
				t.showMessage("Failed to create data. Press Enter to go back.", t.showMainMenu)
				return
			}

			err = t.cache.Set(context.Background(), req.Data.Id, req.Data.Data, 0).Err()
			if err != nil {
				log.Printf("Failed to cache data: %v", err)
			}

			t.showMessage(fmt.Sprintf("Data created successfully.\nID - %s \nPress Enter to go back.", req.Data.Id), t.showMainMenu)
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// getData displays a form for retrieving an existing data item, allowing the user
// to input the ID and type of the data item and sending the get request to the server.
func (t *TUI) getData() {
	form := tview.NewForm()
	form.
		AddInputField("ID", "", 40, nil, nil).
		AddDropDown("Type", []string{"binary", "text", "credentials", "bank card"}, 0, nil).
		AddButton("Submit", func() {
			idField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			_, typeField := form.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()

			if t.client.ServerAvailable {
				req := &proto.GetDataRequest{
					Id:   idField,
					Type: typeField,
				}

				ctx, cancel := t.client.CreateContextWithMetadata(15 * time.Second)
				defer cancel()

				resp, err := t.client.GetData(ctx, req)
				if err != nil {
					t.showMessage("Failed to get data. Press Enter to go back.", t.showMainMenu)
					return
				}
				if len(resp.Data) > 0 {
					if typeField == "binary" {
						fileName := fmt.Sprintf("downloaded_file_%s", idField)
						err = os.WriteFile(fileName, resp.Data[0].Data, 0644)
						if err != nil {
							t.showMessage(fmt.Sprintf("Failed to save file: %v", err), t.showMainMenu)
							return
						}
						t.showMessage(fmt.Sprintf("File downloaded and saved as %s. Press Enter to go back.", fileName), t.showMainMenu)
					} else {
						t.showMessage(formatDataItem(resp.Data[0]), t.showMainMenu)
					}

					err = t.cache.Set(context.Background(), idField, resp.Data[0].Data, 0).Err()
					if err != nil {
						log.Printf("Failed to cache data: %v", err)
					}
				} else {
					t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
				}
			} else {
				data, err := t.cache.Get(context.Background(), idField).Result()
				if err != nil {
					t.showMessage("Failed to get data. Press Enter to go back.", t.showMainMenu)
				}

				t.showMessage(fmt.Sprintf("%s\nServer not available, there is information only about data.", data), t.showMainMenu)
			}
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// listData displays a form for retrieving an existing data items, allowing the user
// to input the ID
func (t *TUI) listData() {
	if !t.client.ServerAvailable {
		t.showMessage("Server not available. Press Enter to go back.", t.showMainMenu)
		return
	}

	ctx, cancel := t.client.CreateContextWithMetadata(15 * time.Second)
	defer cancel()

	resp, err := t.client.ListData(ctx, &emptypb.Empty{})
	if err != nil {
		t.showMessage("Failed to list data. Press Enter to go back.", t.showMainMenu)
		return
	}
	if len(resp.Data) > 0 {
		var builder strings.Builder
		for _, item := range resp.Data {
			err = t.cache.Set(context.Background(), item.Id, item.Data, 0).Err()
			if err != nil {
				log.Printf("Failed to cache data: %v", err)
			}

			builder.WriteString(formatDataItem(item))
		}

		builder.WriteString("Press Enter to go back.")

		t.showMessage(builder.String(), t.showMainMenu)

	} else {
		t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
	}
}

// updateData displays a form for updating an existing data item, allowing the user
// to input the ID, type, data, and metadata, and sending the update request to the server.
func (t *TUI) updateData() {
	if !t.client.ServerAvailable {
		t.showMessage("Server not available. Press Enter to go back.", t.showMainMenu)
		return
	}

	form := tview.NewForm()
	form.
		AddInputField("ID", "", 40, nil, nil).
		AddDropDown("Type", []string{"binary", "text", "credentials", "bank card"}, 0, nil).
		AddInputField("Data", "", 40, nil, nil).
		AddInputField("Meta", "", 20, nil, nil).
		AddButton("Submit", func() {
			idField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			_, typeField := form.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()
			dataField := form.GetFormItemByLabel("Data").(*tview.InputField).GetText()
			metaField := form.GetFormItemByLabel("Meta").(*tview.InputField).GetText()

			req := &proto.UpdateDataRequest{
				Data: &proto.DataItem{
					Id:   idField,
					Type: typeField,
					Data: []byte(dataField),
					Meta: metaField,
				},
			}

			ctx, cancel := t.client.CreateContextWithMetadata(15 * time.Second)
			defer cancel()

			resp, err := t.client.UpdateData(ctx, req)
			if err != nil {
				log.Printf("failed to get data: %v", err)
				t.showMessage("Failed to get data. Press Enter to go back.", t.showMainMenu)
				return
			}
			if len(resp.Message) > 0 {
				log.Printf("UpdateData response: %s", resp.Message)
				t.showMessage("Data updated successfully. Press Enter to go back.", t.showMainMenu)
			} else {
				log.Printf("UpdateData response: no data found")
				t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
			}
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// deleteData displays a form for deleting a data item, allowing the user
// to input the ID and sending the delete request to the server.
func (t *TUI) deleteData() {
	if !t.client.ServerAvailable {
		t.showMessage("Server not available. Press Enter to go back.", t.showMainMenu)
		return
	}

	form := tview.NewForm()
	form.
		AddInputField("ID", "", 40, nil, nil).
		AddButton("Submit", func() {
			idField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()

			req := &proto.DeleteDataRequest{
				Id: idField,
			}

			ctx, cancel := t.client.CreateContextWithMetadata(15 * time.Second)
			defer cancel()

			resp, err := t.client.DeleteData(ctx, req)
			if err != nil {
				t.showMessage("Failed to delete data. Press Enter to go back.", t.showMainMenu)
				return
			}
			if len(resp.Message) > 0 {
				t.showMessage("Data deleted successfully. Press Enter to go back.", t.showMainMenu)

				err = t.cache.Del(context.Background(), idField).Err()
				if err != nil {
					t.showMessage(fmt.Sprintf("Failed to delete data from cache: %v", err), t.showMainMenu)
					return
				}
			} else {
				t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
			}
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// showMessage displays a message to the user with a prompt to press Enter to continue,
// returning to a specified function after the message is acknowledged.
func (t *TUI) showMessage(message string, doneFunc func()) {
	textView := tview.NewTextView().
		SetText(message).
		SetDoneFunc(func(key tcell.Key) {
			doneFunc()
		})
	t.app.SetRoot(textView, true).SetFocus(textView)
}

// restart restarts the TUI application, resetting the interface and returning
// to the initial state.
func (t *TUI) restart() {
	go func() {
		if err := t.Run(); err != nil {
			log.Fatalf("failed to run TUI: %v", err)
		}
	}()
}

// formatDataItem returns specified string format for Data Item.
func formatDataItem(item *proto.DataItem) string {
	return fmt.Sprintf(
		"ID: %s\nType: %s\nData: %s\nMeta: %s\nCreated At: %s\nUpdated At: %s\n\n",
		item.Id, item.Type, string(item.Data), item.Meta,
		item.CreatedAt.AsTime().Format(time.RFC3339),
		item.UpdatedAt.AsTime().Format(time.RFC3339),
	)
}
