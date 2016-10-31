package http

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/member"
)

// MemberCreate creates a new member for the current Org.
func MemberCreate(fn controller.MemberCreateFunc) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentOrg = orgFromContext(ctx)
			p          = payloadMember{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		m, err := fn(currentOrg, p.member)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadMember{member: m})
	}
}

type payloadMember struct {
	member *member.Member
}

type Member struct {
	Email        string    `json:"email"`
	Enabled      bool      `json:"enabled"`
	Firstname    string    `json:"first_name"`
	LastLogin    time.Time `json:"last_login"`
	Lastname     string    `json:"last_name"`
	ID           uint64    `json:"-"`
	OrgID        uint64    `json:"-"`
	Password     string    `json:"password"`
	PublicID     string    `json:"id"`
	PublicOrgID  string    `json:"account_id"`
	SessionToken string    `json:"-"`
	URL          string    `json:"url"`
	Username     string    `json:"user_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *payloadMember) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Email        string    `json:"email"`
		Firstname    string    `json:"first_name"`
		LastLogin    time.Time `json:"last_login"`
		Lastname     string    `json:"last_name"`
		Password     string    `json:"password"`
		PublicID     string    `json:"id"`
		PublicOrgID  string    `json:"account_id"`
		SessionToken string    `json:"-"`
		URL          string    `json:"url"`
		Username     string    `json:"user_name"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{})
}

func (p *payloadMember) UnmarshalJSON(raw []byte) error {
	f := struct {
		Email       string `json:"email"`
		Firstname   string `json:"first_name"`
		Lastname    string `json:"last_name"`
		Password    string `json:"password"`
		PublicID    string `json:"id"`
		PublicOrgID string `json:"account_id"`
		Username    string `json:"user_name"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.member = &member.Member{
		Email:     f.Email,
		Firstname: f.Firstname,
		Lastname:  f.Lastname,
		Password:  f.Password,
		Username:  f.Username,
	}

	return nil
}
