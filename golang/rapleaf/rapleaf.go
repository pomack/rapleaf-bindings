package rapleaf

import (
  "bytes"
  "fmt"
  "http"
  "io/ioutil"
  "net"
  "os"
  "strconv"
  "strings"
  "time"
  "xml"
)

const (
  dateLayout = "2006-01-02"
)

var (
  rapleaf_host = "api.rapleaf.com"
  rapleaf_port = "80"
  ERROR_CODES = map[int]string {
    http.StatusOK : "Request processed successfully.",
    http.StatusAccepted : "This person is currently being searched. Check back shortly and we should have data.",
    http.StatusBadRequest : "Invalid email address or Rapleaf ID.",
    http.StatusUnauthorized : "API key was not provided or is invalid.",
    http.StatusForbidden : "Your query limit has been exceeded. Contact developer@rapleaf.com if you would like to increase your limit.",
    http.StatusNotFound : "Returned for lookup by hash or site userid. We do not have this person in our system. If you would like better results, consider supplying the email address.",
    http.StatusInternalServerError : "There was an unexpected error on our server. This should be very rare and if you see it please contact developer@rapleaf.com.",
  }
)

func OverrideRapleafHostPort(host, port string) {
  rapleaf_host = host
  rapleaf_port = port
}

func personUrl(a ...string) string {
  arr := []string{"http://", rapleaf_host, ":", rapleaf_port, "/v3/person/"}
  v := make([]string, len(a) + len(arr))
  for i, e := range arr {
    v[i] = e
  }
  for i, e := range a {
    v[i + len(arr)] = e
  }
  return strings.Join(v, "");
}

func graphUrl(a ...string) string {
  arr := []string{"http://", rapleaf_host, ":", rapleaf_port, "/v2/graph/"}
  v := make([]string, len(a) + len(arr))
  for i, e := range arr {
    v[i] = e
  }
  for i, e := range a {
    v[i + len(arr)] = e
  }
  return strings.Join(v, "");
}

type rapleafMemberSite struct {
  XMLName xml.Name "membership"
  Site string "attr"
  Exists string "attr"
  Profile_url string "attr"
  Image_url string "attr"
  Num_friends string "attr"
  Num_followers string "attr"
  Num_followed string "attr"
}

type rapleafOccupation struct {
  XMLName xml.Name "occupation"
  Company string "attr"
  Job_title string "attr"
}

type rapleafPrimaryMembership struct {
  XMLName xml.Name "primary"
  Membership []rapleafMemberSite
}

type rapleafSupplementalMembership struct {
  XMLName xml.Name "supplemental"
  Membership []rapleafMemberSite
}

type rapleafMemberships struct {
  XMLName xml.Name "memberships"
  Primary rapleafPrimaryMembership
  Supplemental rapleafSupplementalMembership
}

type rapleafOccupations struct {
  XMLName xml.Name "occupations"
  Occupation []rapleafOccupation
}

type rapleafBasics struct {
  XMLName xml.Name "basics"
  Name string
  Gender string
  Location string
  Num_friends int
  Age int
  Earliest_known_activity string
  Latest_known_activity string
  Occupations []rapleafOccupations
}

type rapleafPerson struct {
  XMLName xml.Name "person"
  Id string "attr"
  Basics rapleafBasics
  Memberships rapleafMemberships
}

type RapleafMemberSite struct {
  Site string
  ProfileUrl string
  ImageUrl string
  NumFriends int
  NumFollowers int
  NumFollowed int
  Exists string
}

type RapleafOccupation struct {
  Company string
  JobTitle string
}

type RapleafPerson struct {
  Id string
  Name string
  Gender string
  Location string
  NumFriends int
  Age int
  EarliestKnownActivity *time.Time
  LatestKnownActivity *time.Time
  Occupations []*RapleafOccupation
  Memberships []*RapleafMemberSite
  EmailAddress string
}

func (p* rapleafMemberSite) toPublicStruct() *RapleafMemberSite {
  friends, _ := strconv.Atoi(p.Num_friends)
  followers, _ := strconv.Atoi(p.Num_followers)
  followed, _ := strconv.Atoi(p.Num_followed)
  return &RapleafMemberSite{
    Site:p.Site,
    ProfileUrl:p.Profile_url,
    ImageUrl:p.Image_url,
    NumFriends:friends,
    NumFollowers:followers,
    NumFollowed:followed,
    Exists:p.Exists,
  }
}

func (p* rapleafOccupation) toPublicStruct() *RapleafOccupation {
  return &RapleafOccupation{
    Company:p.Company,
    JobTitle:p.Job_title,
  }
}

func (p* rapleafPerson) toPublicStruct() *RapleafPerson {
  earliest_known_activity, _ := time.Parse(dateLayout, p.Basics.Earliest_known_activity)
  latest_known_activity, _ := time.Parse(dateLayout, p.Basics.Latest_known_activity)
  num_occupations := 0
  if len(p.Basics.Occupations) > 0 {
    num_occupations = len(p.Basics.Occupations[0].Occupation)
  }
  occupations := make([]*RapleafOccupation, num_occupations)
  memberships := make([]*RapleafMemberSite, len(p.Memberships.Primary.Membership) + len(p.Memberships.Supplemental.Membership))
  if num_occupations > 0 {
    for j, occupation := range p.Basics.Occupations[0].Occupation {
      occupations[j] = occupation.toPublicStruct()
    }
  }
  i := 0
  for _, membership := range p.Memberships.Primary.Membership {
    memberships[i] = membership.toPublicStruct()
    i++
  }
  for _, membership := range p.Memberships.Supplemental.Membership {
    memberships[i] = membership.toPublicStruct()
    i++
  }
  return &RapleafPerson{
    Id:p.Id,
    Name:p.Basics.Name,
    Gender:strings.ToLower(p.Basics.Gender),
    Location:p.Basics.Location,
    NumFriends:p.Basics.Num_friends,
    Age:p.Basics.Age,
    EarliestKnownActivity:earliest_known_activity,
    LatestKnownActivity:latest_known_activity,
    Occupations:occupations,
    Memberships:memberships,
  }
}

func (p *RapleafMemberSite) String() string {
  arr := []string{"rapleaf.RapleafMemberSite{",
    "Site:", strconv.Quote(p.Site), ", ",
    "ProfileUrl:", strconv.Quote(p.ProfileUrl), ", ",
    "ImageUrl:", strconv.Quote(p.ImageUrl), ", ",
    "NumFriends:", strconv.Itoa(p.NumFriends), ", ",
    "NumFollowers:", strconv.Itoa(p.NumFollowers), ", ",
    "NumFollowed:", strconv.Itoa(p.NumFollowed), ", ",
    "Exists:", strconv.Quote(p.Exists), "}",
  }
  return strings.Join(arr, "")
}

func (p *RapleafMemberSite) Equals(other *RapleafMemberSite) bool {
  return other != nil && p.ProfileUrl == other.ProfileUrl && 
    p.ImageUrl == other.ImageUrl && 
    p.NumFriends == other.NumFriends &&
    p.NumFollowers == other.NumFollowers &&
    p.NumFollowed == other.NumFollowed &&
    p.Exists == other.Exists
}

func (p *RapleafOccupation) String() string {
  arr := []string{"rapleaf.RapleafOccupation{",
    "Company:", strconv.Quote(p.Company), ", ",
    "JobTitle:", strconv.Quote(p.JobTitle), "}",
  }
  return strings.Join(arr, "")
}

func (p *RapleafOccupation) Equals(other *RapleafOccupation) bool {
  return other != nil && p.Company == other.Company && 
    p.JobTitle == other.JobTitle
}

func (p *RapleafPerson) String() string {
  var earliest_known_activity_str string
  var latest_known_activity_str string
  if p.EarliestKnownActivity != nil && p.EarliestKnownActivity.Year > 1000 {
    earliest_known_activity_str = fmt.Sprintf("&time.Time{Year:%d, Month:%d, Day:%d}", p.EarliestKnownActivity.Year, p.EarliestKnownActivity.Month, p.EarliestKnownActivity.Day)
  } else {
    earliest_known_activity_str = "nil"
  }
  if p.LatestKnownActivity != nil && p.LatestKnownActivity.Year > 1000 {
    latest_known_activity_str = fmt.Sprintf("&time.Time{Year:%d, Month:%d, Day:%d}", p.LatestKnownActivity.Year, p.LatestKnownActivity.Month, p.LatestKnownActivity.Day)
  } else {
    latest_known_activity_str = "nil"
  }
  occupations := make([]string, len(p.Occupations))
  memberships := make([]string, len(p.Memberships))
  for i, occupation := range p.Occupations {
    occupations[i] = "&" + occupation.String()
  }
  for i, membership := range p.Memberships {
    memberships[i] = "&" + membership.String()
  }
  arr := []string{"rapleaf.RapleafPerson{",
    "Id:", strconv.Quote(p.Id), ", ",
    "Name:", strconv.Quote(p.Name), ", ",
    "Gender:", strconv.Quote(p.Gender), ", ",
    "Location:", strconv.Quote(p.Location), ", ",
    "NumFriends:", strconv.Itoa(p.NumFriends), ", ",
    "Age:", strconv.Itoa(p.Age), ", ",
    "EarliestKnownActivity:", earliest_known_activity_str, ", ",
    "LatestKnownActivity:", latest_known_activity_str, ", ",
    "EmailAddress:", strconv.Quote(p.EmailAddress), ", ",
    "Occupations:[]*rapleaf.RapleafOccupation{", strings.Join(occupations, ", "), "}, ",
    "Memberships:[]*rapleaf.RapleafMemberSite{", strings.Join(memberships, ", "), "} }",
  }
  return strings.Join(arr, "")
}

func (p *RapleafPerson) Equals(other *RapleafPerson) bool {
  if other == nil {
    return false
  }
  if p.Id != other.Id ||
      p.Name != other.Name ||
      p.Gender != other.Gender ||
      p.Location != other.Location ||
      p.NumFriends != other.NumFriends ||
      p.Age != other.Age ||
      p.EmailAddress != other.EmailAddress ||
      len(p.Occupations) != len(other.Occupations) ||
      len(p.Memberships) != len(other.Memberships) {
    return false
  }
  if p.EarliestKnownActivity != other.EarliestKnownActivity {
    if(p.EarliestKnownActivity == nil || other.EarliestKnownActivity == nil) {
      return false
    }
    if p.EarliestKnownActivity.Seconds() != other.EarliestKnownActivity.Seconds() {
      return false
    }
  }
  if p.LatestKnownActivity != other.LatestKnownActivity {
    if(p.LatestKnownActivity == nil || other.LatestKnownActivity == nil) {
      return false
    }
    if p.LatestKnownActivity.Seconds() != other.LatestKnownActivity.Seconds() {
      return false
    }
  }
  for i, occupation := range p.Occupations {
    if !occupation.Equals(other.Occupations[i]) {
      return false
    }
  }
  for i, membership := range p.Memberships {
    if !membership.Equals(other.Memberships[i]) {
      return false
    }
  }
  return true
}

func RapleafPersonFromString(value string) (*RapleafPerson, os.Error) {
  if(len(value) == 0) { return nil, nil; }
  b := bytes.NewBufferString(value)
  p := &rapleafPerson{}
  if err := xml.Unmarshal(b, p); err != nil {
    return nil, err
  }
  return p.toPublicStruct(), nil
}

func retrieve(api_key, url string) (code int, text string) {
  parsedUrl, err := http.ParseURL(url)
  if err != nil {
    return http.StatusBadRequest, err.String()
  }
  headers := make(map[string]string)
  headers["Authorization"] = api_key
  req := &http.Request{Method:"GET", 
    RawURL:url, 
    URL: parsedUrl, 
    Proto:"HTTP/1.1", 
    ProtoMajor:1, 
    ProtoMinor:1, 
    Header:headers, 
    Host:rapleaf_host,
  }
  c, err := net.Dial("tcp", "", rapleaf_host + ":" + rapleaf_port)
  if err != nil {
    return http.StatusServiceUnavailable, err.String()
  }
  conn := http.NewClientConn(c, nil)
  if err := conn.Write(req); err != nil {
    return http.StatusServiceUnavailable, err.String()
  }
  resp, err := conn.Read()
  if resp == nil {
    if err == nil {
      return http.StatusServiceUnavailable, err.String()
    }
    return http.StatusNoContent, ""
  }
  buf, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return resp.StatusCode, err.String()
  }
  return resp.StatusCode, string(buf)
}

func PersonXmlByEmail(api_key, email_address string) (int, string) {
  url := personUrl("email/", http.URLEscape(email_address))
  return retrieve(api_key, url)
}

func PersonXmlByRapleafId(api_key, rapleaf_id string) (int, string) {
  return PersonXmlBySite(api_key, "rapleaf", rapleaf_id)
}

func PersonXmlBySite(api_key, site, profile_id string) (int, string) {
  url := personUrl("web/", http.URLEscape(site), "/", http.URLEscape(profile_id))
  return retrieve(api_key, url)
}

func PersonByEmail(api_key, email_address string) (*RapleafPerson) {
  code, text := PersonXmlByEmail(api_key, email_address)
  if code == http.StatusOK {
    u, err := RapleafPersonFromString(text)
    if err == nil && u != nil {
      u.EmailAddress = email_address
      return u
    }
    return u
  }
  return nil
}

func PersonByRapleafId(api_key, rapleaf_id string) (*RapleafPerson) {
  return PersonBySite(api_key, "rapleaf", rapleaf_id)
}

func PersonBySite(api_key, site, profile_id string) (*RapleafPerson) {
  code, text := PersonXmlBySite(api_key, site, profile_id)
  if code == http.StatusOK {
    u, err := RapleafPersonFromString(text)
    if err == nil {
      return u
    }
  }
  return nil
}

