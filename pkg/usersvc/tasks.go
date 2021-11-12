package usersvc

type ConfirmationEmailTask struct {
	UserID string
}

func (t *ConfirmationEmailTask) GetTaskID() string {
	return t.UserID
}

type PasswordResetEmailTask struct {
	Email string
}

func (t *PasswordResetEmailTask) GetTaskID() string {
	return t.Email
}

type emailUpdateConfirmationTask struct {
	UserID   string
	NewEmail string
}

func (t *emailUpdateConfirmationTask) GetTaskID() string {
	return t.UserID + "-" + t.NewEmail
}
func (usvc *UserService) registerTaskTypes() {
	usvc.svc.JQ.Register(&ConfirmationEmailTask{}, usvc)
	usvc.svc.JQ.Register(&PasswordResetEmailTask{}, usvc)
	usvc.svc.JQ.Register(&emailUpdateConfirmationTask{}, usvc)
}
func (usvc *UserService) Execute(i interface{}) error {

	switch v := i.(type) {
	case *ConfirmationEmailTask:
		usvc.sendConfirmationEmail(v.UserID)
	case *PasswordResetEmailTask:
		usvc.sendPasswordResetEmail(v.Email)
	case *emailUpdateConfirmationTask:
		usvc.sendEmailUpdateEmail(v.UserID, v.NewEmail)
	}

	return nil
}
