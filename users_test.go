package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"
)

func getTestUserProfileCustomField() UserProfileCustomField {
	return UserProfileCustomField{
		Value: "test value",
		Alt:   "",
		Label: "",
	}
}

func getTestUserProfileCustomFields() UserProfileCustomFields {
	return UserProfileCustomFields{
		fields: map[string]UserProfileCustomField{
			"Xxxxxx": getTestUserProfileCustomField(),
		}}
}

func getTestUserProfile() UserProfile {
	return UserProfile{
		StatusText:            "testStatus",
		StatusEmoji:           ":construction:",
		RealName:              "Test Real Name",
		RealNameNormalized:    "Test Real Name Normalized",
		DisplayName:           "Test Display Name",
		DisplayNameNormalized: "Test Display Name Normalized",
		Email:    "test@test.com",
		Image24:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_24.jpg",
		Image32:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_32.jpg",
		Image48:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_48.jpg",
		Image72:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_72.jpg",
		Image192: "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_192.jpg",
		Fields:   getTestUserProfileCustomFields(),
	}
}

func getTestUser() User {
	return User{
		ID:                "UXXXXXXXX",
		Name:              "Test User",
		Deleted:           false,
		Color:             "9f69e7",
		RealName:          "testuser",
		TZ:                "America/Los_Angeles",
		TZLabel:           "Pacific Daylight Time",
		TZOffset:          -25200,
		Profile:           getTestUserProfile(),
		IsBot:             false,
		IsAdmin:           false,
		IsOwner:           false,
		IsPrimaryOwner:    false,
		IsRestricted:      false,
		IsUltraRestricted: false,
		Has2FA:            false,
	}
}

func getUserIdentity(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
  "ok": true,
  "user": {
    "id": "UXXXXXXXX",
    "name": "Test User",
    "email": "test@test.com",
    "image_24": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_24.jpg",
    "image_32": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_32.jpg",
    "image_48": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_48.jpg",
    "image_72": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_72.jpg",
    "image_192": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_192.jpg",
    "image_512": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_512.jpg"
  },
  "team": {
    "id": "TXXXXXXXX",
    "name": "team-name",
    "domain": "team-domain",
    "image_34": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_34.jpg",
    "image_44": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_44.jpg",
    "image_68": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_68.jpg",
    "image_88": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_88.jpg",
    "image_102": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_102.jpg",
    "image_132": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_132.jpg",
    "image_230": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_230.jpg",
    "image_original": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_original.jpg"
  }
}`)
	rw.Write(response)
}

func getUserByEmail(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		Ok   bool `json:"ok"`
		User User `json:"user"`
	}{
		Ok:   true,
		User: getTestUser(),
	})
	rw.Write(response)
}

func httpTestErrReply(w http.ResponseWriter, clientErr bool, msg string) {
	if clientErr {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	body, _ := json.Marshal(&SlackResponse{
		Ok: false, Error: msg,
	})

	w.Write(body)
}

func newProfileHandler(up *UserProfile) (setter func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		if up == nil {
			httpTestErrReply(w, false, "err: UserProfile is nil")
			return
		}

		if err := r.ParseForm(); err != nil {
			httpTestErrReply(w, true, fmt.Sprintf("err parsing form: %s", err.Error()))
			return
		}

		values := r.Form

		if len(values["profile"]) == 0 {
			httpTestErrReply(w, true, `POST data must include a "profile" field`)
			return
		}

		profile := []byte(values["profile"][0])

		userProfile := UserProfile{}

		if err := json.Unmarshal(profile, &userProfile); err != nil {
			httpTestErrReply(w, true, fmt.Sprintf("err parsing JSON: %s\n\njson: `%s`", err.Error(), profile))
			return
		}

		*up = userProfile

		// TODO(theckman): enhance this to return a full User object
		fmt.Fprint(w, `{"ok":true}`)
	}
}

func TestGetUserIdentity(t *testing.T) {
	http.HandleFunc("/users.identity", getUserIdentity)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	identity, err := api.GetUserIdentity()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	// t.Fatal refers to -> t.Errorf & return
	if identity.User.ID != "UXXXXXXXX" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.User.Name != "Test User" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.User.Email != "test@test.com" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.ID != "TXXXXXXXX" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.Name != "team-name" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.Domain != "team-domain" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.User.Image24 == "" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.Image34 == "" {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetUserByEmail(t *testing.T) {
	http.HandleFunc("/users.lookupByEmail", getUserByEmail)
	expectedUser := getTestUser()

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	user, err := api.GetUserByEmail("test@test.com")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedUser, *user) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestUserCustomStatus(t *testing.T) {
	up := &UserProfile{}

	setUserProfile := newProfileHandler(up)

	http.HandleFunc("/users.profile.set", setUserProfile)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	testSetUserCustomStatus(api, up, t)
	testUnsetUserCustomStatus(api, up, t)
}

func testSetUserCustomStatus(api *Client, up *UserProfile, t *testing.T) {
	const (
		statusText  = "testStatus"
		statusEmoji = ":construction:"
	)

	if err := api.SetUserCustomStatus(statusText, statusEmoji); err != nil {
		t.Fatalf(`SetUserCustomStatus(%q, %q) = %#v, want <nil>`, statusText, statusEmoji, err)
	}

	if up.StatusText != statusText {
		t.Fatalf(`UserProfile.StatusText = %q, want %q`, up.StatusText, statusText)
	}

	if up.StatusEmoji != statusEmoji {
		t.Fatalf(`UserProfile.StatusEmoji = %q, want %q`, up.StatusEmoji, statusEmoji)
	}
}

func testUnsetUserCustomStatus(api *Client, up *UserProfile, t *testing.T) {
	if err := api.UnsetUserCustomStatus(); err != nil {
		t.Fatalf(`UnsetUserCustomStatus() = %#v, want <nil>`, err)
	}

	if up.StatusText != "" {
		t.Fatalf(`UserProfile.StatusText = %q, want %q`, up.StatusText, "")
	}

	if up.StatusEmoji != "" {
		t.Fatalf(`UserProfile.StatusEmoji = %q, want %q`, up.StatusEmoji, "")
	}
}

func TestGetUsers(t *testing.T) {
	http.HandleFunc("/users.list", getUserPage(4))
	expectedUser := getTestUser()

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	users, err := api.GetUsers()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if !reflect.DeepEqual([]User{expectedUser, expectedUser, expectedUser, expectedUser}, users) {
		t.Fatal(ErrIncorrectResponse)
	}
}

// returns n pages users.
func getUserPage(max int64) func(rw http.ResponseWriter, r *http.Request) {
	var n int64
	return func(rw http.ResponseWriter, r *http.Request) {
		var cpage int64
		sresp := SlackResponse{
			Ok: true,
		}
		members := []User{
			getTestUser(),
		}
		rw.Header().Set("Content-Type", "application/json")
		if cpage = atomic.AddInt64(&n, 1); cpage == max {
			response, _ := json.Marshal(userResponseFull{
				SlackResponse: sresp,
				Members:       members,
			})
			rw.Write(response)
			return
		}
		response, _ := json.Marshal(userResponseFull{
			SlackResponse: sresp,
			Members:       members,
			Metadata:      ResponseMetadata{Cursor: strconv.Itoa(int(cpage))},
		})
		rw.Write(response)
	}
}

func TestSetUserPhoto(t *testing.T) {
	file, fileContent, teardown := createUserPhoto(t)
	defer teardown()

	params := UserSetPhotoParams{CropX: 0, CropY: 0, CropW: 32}

	http.HandleFunc("/users.setPhoto", setUserPhotoHandler(fileContent, params))

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New(validToken)

	err := api.SetUserPhoto(file.Name(), params)
	if err != nil {
		t.Fatalf("unexpected error: %+v\n", err)
	}
}

func setUserPhotoHandler(wantBytes []byte, wantParams UserSetPhotoParams) http.HandlerFunc {
	const maxMemory = 1 << 20 // 1 MB

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			httpTestErrReply(w, false, fmt.Sprintf("failed to parse multipart/form: %+v", err))
			return
		}

		// Test for expected token
		if v := r.Form.Get("token"); v != validToken {
			httpTestErrReply(w, true, fmt.Sprintf("expected multipart form value token=%v", validToken))
			return
		}

		// Test for expected crop params
		if wantParams.CropX != DEFAULT_USER_PHOTO_CROP_X {
			if cx, err := strconv.Atoi(r.Form.Get("crop_x")); err != nil || cx != wantParams.CropX {
				httpTestErrReply(w, true, fmt.Sprintf("expected multipart form value crop_x=%d", wantParams.CropX))
				return
			}
		}
		if wantParams.CropY != DEFAULT_USER_PHOTO_CROP_Y {
			if cy, err := strconv.Atoi(r.Form.Get("crop_y")); err != nil || cy != wantParams.CropY {
				httpTestErrReply(w, true, fmt.Sprintf("expected multipart form value crop_y=%d", wantParams.CropY))
				return
			}
		}
		if wantParams.CropW != DEFAULT_USER_PHOTO_CROP_W {
			if cw, err := strconv.Atoi(r.Form.Get("crop_w")); err != nil || cw != wantParams.CropW {
				httpTestErrReply(w, true, fmt.Sprintf("expected multipart form value crop_w=%d", wantParams.CropW))
				return
			}
		}

		// Test for expected image
		f, ok := r.MultipartForm.File["image"]
		if !ok || len(f) == 0 {
			httpTestErrReply(w, true, `expected multipart form file "image"`)
			return
		}
		file, err := f[0].Open()
		if err != nil {
			httpTestErrReply(w, true, fmt.Sprintf("failed to open uploaded file: %+v", err))
			return
		}
		gotBytes, err := ioutil.ReadAll(file)
		if err != nil {
			httpTestErrReply(w, true, fmt.Sprintf("failed to read uploaded file: %+v", err))
			return
		}
		if !bytes.Equal(wantBytes, gotBytes) {
			httpTestErrReply(w, true, "uploaded bytes did not match expected bytes")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true}`)
	}
}

// createUserPhoto generates a temp photo for testing. It returns the file handle, the file
// contents, and a function that can be called to remove the file.
func createUserPhoto(t *testing.T) (*os.File, []byte, func()) {
	photo := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.Draw(photo, photo.Bounds(), image.Black, image.ZP, draw.Src)

	f, err := ioutil.TempFile(os.TempDir(), "profile.png")
	if err != nil {
		t.Fatalf("failed to create test photo: %+v\n", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(io.MultiWriter(&buf, f), photo); err != nil {
		t.Fatalf("failed to write test photo: %+v\n", err)
	}

	teardown := func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Fatalf("failed to remove test photo: %+v\n", err)
		}
	}

	return f, buf.Bytes(), teardown
}

func getUserProfileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	profile := getTestUserProfile()
	resp, _ := json.Marshal(&getUserProfileResponse{
		SlackResponse: SlackResponse{Ok: true},
		Profile:       &profile})
	rw.Write(resp)
}

func TestGetUserProfile(t *testing.T) {
	http.HandleFunc("/users.profile.get", getUserProfileHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	profile, err := api.GetUserProfile("UXXXXXXXX", false)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	exp := getTestUserProfile()
	if profile.DisplayName != exp.DisplayName {
		t.Fatalf(`profile.DisplayName = "%s", wanted "%s"`, profile.DisplayName, exp.DisplayName)
	}
}

func TestSetFieldsMap(t *testing.T) {
	p := &UserProfile{}
	exp := map[string]UserProfileCustomField{
		"Xxxxxx": getTestUserProfileCustomField(),
	}
	p.SetFieldsMap(exp)
	act := p.FieldsMap()
	if !reflect.DeepEqual(act, exp) {
		t.Fatalf(`p.FieldsMap() = %v, wanted %v`, act, exp)
	}
}

func TestUserProfileCustomFieldsUnmarshalJSON(t *testing.T) {
	fields := &UserProfileCustomFields{}
	if err := json.Unmarshal([]byte(`[]`), fields); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`{
	  "Xxxxxx": {
	    "value": "test value",
	    "alt": ""
	  }
	}`), fields); err != nil {
		t.Fatal(err)
	}
	act := fields.ToMap()["Xxxxxx"].Value
	exp := "test value"
	if act != exp {
		t.Fatalf(`fields.ToMap()["Xxxxxx"]["value"] = "%s", wanted "%s"`, act, exp)
	}
}

func TestUserProfileCustomFieldsMarshalJSON(t *testing.T) {
	fields := UserProfileCustomFields{}
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "[]" {
		t.Fatalf(`string(b) = "%s", wanted "[]"`, string(b))
	}
	fields = getTestUserProfileCustomFields()
	if _, err := json.Marshal(fields); err != nil {
		t.Fatal(err)
	}
}

func TestUserProfileCustomFieldsToMap(t *testing.T) {
	m := map[string]UserProfileCustomField{
		"Xxxxxx": getTestUserProfileCustomField(),
	}
	fields := UserProfileCustomFields{fields: m}
	act := fields.ToMap()
	if !reflect.DeepEqual(act, m) {
		t.Fatalf(`fields.ToMap() = %v, wanted %v`, act, m)
	}
}

func TestUserProfileCustomFieldsLen(t *testing.T) {
	fields := UserProfileCustomFields{
		fields: map[string]UserProfileCustomField{
			"Xxxxxx": getTestUserProfileCustomField(),
		}}
	if fields.Len() != 1 {
		t.Fatalf(`fields.Len() = %d, wanted 1`, fields.Len())
	}
}

func TestUserProfileCustomFieldsSetMap(t *testing.T) {
	fields := UserProfileCustomFields{}
	m := map[string]UserProfileCustomField{
		"Xxxxxx": getTestUserProfileCustomField(),
	}
	fields.SetMap(m)
	if !reflect.DeepEqual(fields.fields, m) {
		t.Fatalf(`fields.fields = %v, wanted %v`, fields.fields, m)
	}
}
