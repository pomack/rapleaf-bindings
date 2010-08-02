package rapleaf_test

/*
 * Copyright 2010 Aalok Shah (aalok@shah.ws)
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *      http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
  . "rapleaf"
  "http"
  "net"
  "os"
  "strconv"
  "strings"
  "testing"
  "time"
  "rand"
)

const (
  API_KEY = "stuff"
  USER_EMPTY_XML = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><person id=\"b34282025d7e2c5db6786a8daaab48c7\"><basics><earliest_known_activity>2010-05-27</earliest_known_activity><num_friends>0</num_friends></basics><memberships><primary><membership site=\"bebo.com\" exists=\"false\"/><membership site=\"facebook.com\" exists=\"unknown\"/><membership site=\"flickr.com\" exists=\"false\"/><membership site=\"friendster.com\" exists=\"false\"/><membership site=\"hi5.com\" exists=\"false\"/><membership site=\"linkedin.com\" exists=\"tbd\"/><membership site=\"livejournal.com\" exists=\"false\"/><membership site=\"metroflog.com\" exists=\"false\"/><membership site=\"multiply.com\" exists=\"unknown\"/><membership site=\"myspace.com\" exists=\"false\"/><membership site=\"myyearbook.com\" exists=\"false\"/><membership site=\"plaxo.com\" exists=\"false\"/><membership site=\"twitter.com\" exists=\"unknown\"/></primary><supplemental></supplemental></memberships></person>"
  USER_WITH_PROFILE_XML = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><person id=\"97fc425100000000\"><basics><name>John Q Public</name><age>28</age><gender>Male</gender><location>Albuquerque, New Mexico, United States</location><occupations><occupation job_title=\"Software Developer\" company=\"Apple\" /><occupation job_title=\"VP Marketing\" company=\"GE\" /><occupation job_title=\"Founder\" company=\"Startup.com\" /></occupations><earliest_known_activity>2001-11-16</earliest_known_activity><latest_known_activity>2010-05-08</latest_known_activity><num_friends>156</num_friends></basics><memberships><primary><membership site=\"bebo.com\" exists=\"false\"/><membership site=\"facebook.com\" exists=\"true\"/><membership site=\"flickr.com\" exists=\"false\"/><membership site=\"friendster.com\" exists=\"true\" profile_url=\"http://profiles.friendster.com/3543228\" image_url=\"http://photos.friendster.com/photos/82/11/3543228/13281738852124s.jpg\" num_friends=\"16\"/><membership site=\"hi5.com\" exists=\"false\"/><membership site=\"linkedin.com\" exists=\"true\" profile_url=\"http://www.linkedin.com/in/johnqpublic\" image_url=\"http://media.linkedin.com/mpr/mpr/shrink_80_80/p/2/000/016/0f0/36426ef.jpg\" num_friends=\"166\"/><membership site=\"livejournal.com\" exists=\"false\"/><membership site=\"metroflog.com\" exists=\"false\"/><membership site=\"multiply.com\" exists=\"false\"/><membership site=\"myspace.com\" exists=\"false\"/><membership site=\"myyearbook.com\" exists=\"false\"/><membership site=\"plaxo.com\" exists=\"false\"/><membership site=\"twitter.com\" exists=\"true\" profile_url=\"http://twitter.com/johnqpublic\" num_followers=\"14\" num_followed=\"4\"/></primary><supplemental><membership site=\"pandora.com\" exists=\"true\" profile_url=\"http://www.pandora.com/people/johnqpublic\"/><membership site=\"tagged.com\" exists=\"true\" profile_url=\"http://www.tagged.com/profile.html?uid=5378192615\" num_friends=\"0\" num_followers=\"0\" num_followed=\"0\"/></supplemental></memberships></person>"
)

var (
  URL_MAPPINGS = map[string]string {
    "/v3/person/email/empty.profile@gmail.com" : USER_EMPTY_XML,
    "/v3/person/web/rapleaf/b34282025d7e2c5db6786a8daaab48c7" : USER_EMPTY_XML,
    "/v3/person/email/john.q.public@gmail.com" : USER_WITH_PROFILE_XML,
    "/v3/person/web/friendster/3543228" : USER_WITH_PROFILE_XML,
    "/v3/person/web/linkedin/johnqpublic" : USER_WITH_PROFILE_XML,
    "/v3/person/web/pandora/johnqpublic" : USER_WITH_PROFILE_XML,
    "/v3/person/web/rapleaf/97fc425100000000" : USER_WITH_PROFILE_XML,
    "/v3/person/web/tagged/5378192615" : USER_WITH_PROFILE_XML,
    "/v3/person/web/twitter/johnqpublic" : USER_WITH_PROFILE_XML,
  }
)

var (
  USER_EMPTY_PERSON = &RapleafPerson{
    Id:"b34282025d7e2c5db6786a8daaab48c7", 
    EarliestKnownActivity:&time.Time{Year: 2010, Month: 5, Day: 27},
    Memberships:[]*RapleafMemberSite{
      &RapleafMemberSite{
        Site:"bebo.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"facebook.com",
        Exists:"unknown",
      },
      &RapleafMemberSite{
        Site:"flickr.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"friendster.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"hi5.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"linkedin.com",
        Exists:"tbd",
      },
      &RapleafMemberSite{
        Site:"livejournal.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"metroflog.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"multiply.com",
        Exists:"unknown",
      },
      &RapleafMemberSite{
        Site:"myspace.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"myyearbook.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"plaxo.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"twitter.com",
        Exists:"unknown",
      },
    },
  }
  
  USER_WITH_PROFILE_PERSON = &RapleafPerson{
    Id:"97fc425100000000",
    Name:"John Q Public",
    Gender:"male",
    Location:"Albuquerque, New Mexico, United States",
    NumFriends:156,
    Age:28,
    EarliestKnownActivity:&time.Time{Year:2001, Month:11, Day:16},
    LatestKnownActivity:&time.Time{Year:2010, Month:5, Day:8},
    Occupations:[]*RapleafOccupation{
      &RapleafOccupation{
        Company:"Apple",
        JobTitle:"Software Developer",
      },
      &RapleafOccupation{
        Company:"GE",
        JobTitle:"VP Marketing",
      },
      &RapleafOccupation{
        Company:"Startup.com", 
        JobTitle:"Founder",
      },
    },
    Memberships:[]*RapleafMemberSite{
      &RapleafMemberSite{
        Site:"bebo.com", 
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"facebook.com", 
        Exists:"true",
      },
      &RapleafMemberSite{
        Site:"flickr.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"friendster.com",
        ProfileUrl:"http://profiles.friendster.com/3543228",
        ImageUrl:"http://photos.friendster.com/photos/82/11/3543228/13281738852124s.jpg",
        NumFriends:16,
        Exists:"true",
      },
      &RapleafMemberSite{
        Site:"hi5.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"linkedin.com",
        ProfileUrl:"http://www.linkedin.com/in/johnqpublic",
        ImageUrl:"http://media.linkedin.com/mpr/mpr/shrink_80_80/p/2/000/016/0f0/36426ef.jpg",
        NumFriends:166,
        Exists:"true",
      },
      &RapleafMemberSite{
        Site:"livejournal.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"metroflog.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"multiply.com",
         Exists:"false",
      },
      &RapleafMemberSite{
        Site:"myspace.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"myyearbook.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"plaxo.com",
        Exists:"false",
      },
      &RapleafMemberSite{
        Site:"twitter.com",
        ProfileUrl:"http://twitter.com/johnqpublic",
        NumFollowers:14,
        NumFollowed:4,
        Exists:"true",
      },
      &RapleafMemberSite{
        Site:"pandora.com",
        ProfileUrl:"http://www.pandora.com/people/johnqpublic",
        Exists:"true",
      },
      &RapleafMemberSite{
        Site:"tagged.com",
        ProfileUrl:"http://www.tagged.com/profile.html?uid=5378192615",
        Exists:"true",
      },
    },
  }
)

func ServeTestHTTP(conn *http.Conn, req *http.Request) {
  req.Close = true
  if api_key, ok := req.Header["Authorization"]; !ok || api_key != API_KEY {
    text := []byte(ERROR_CODES[http.StatusUnauthorized])
    conn.SetHeader("Content-Type", "text/html;charset=ISO-8859-1")
    conn.SetHeader("Cache-Control", "must-revalidate,no-cache,no-store")
    conn.SetHeader("Content-Length", strconv.Itoa(len(text)))
    conn.WriteHeader(http.StatusUnauthorized)
    conn.Write(text)
    conn.Flush()
    return
  }
  if text, ok := URL_MAPPINGS[req.URL.Path]; ok {
    conn.SetHeader("Content-Type", "application/xml;charset=UTF-8")
    conn.SetHeader("Transfer-Encoding", "chunked")
    // for some reason, api.rapleaf.com does not send Content-Length
    // when sending stored data
    conn.WriteHeader(http.StatusOK)
    conn.Write([]byte(text))
    conn.Flush()
    return
  }
  text := []byte(ERROR_CODES[http.StatusNotFound])
  conn.SetHeader("Content-Type", "text/html;charset=ISO-8859-1")
  conn.SetHeader("Cache-Control", "must-revalidate,no-cache,no-store")
  conn.SetHeader("Content-Length", strconv.Itoa(len(text)))
  conn.WriteHeader(http.StatusNotFound)
  conn.Write(text)
  conn.Flush()
  return
}

func serveTestFiles(t *testing.T) (l net.Listener, err os.Error) {
  foundValidPort := false
  var port_str string
  for !foundValidPort {
    port := (rand.Int() & 0x7FFF) + 0x08000
    port_str = strconv.Itoa(port)
    addr, err := net.ResolveTCPAddr("127.0.0.1:" + port_str)
    if err != nil {
      t.Error("Create TCP Address: ", err.String())
      return nil, err
    }
    l, err = net.ListenTCP("tcp4", addr)
    if err != nil {
      if err == os.EADDRINUSE || strings.LastIndex(err.String(), os.EADDRINUSE.String()) != -1 {
        continue
      }
      t.Error("Unable to listen on TCP port: ", err.String())
      return l, err
    }
    foundValidPort = true
  }
  OverrideRapleafHostPort("127.0.0.1", port_str)
  go http.Serve(l, http.HandlerFunc(ServeTestHTTP))
  return l, err
}

func closeServerTestFiles(l net.Listener) {
  if l != nil {
    l.Close()
  }
}

func TestSetup(t *testing.T) {
  serveTestFiles(t)
}

func testSameOccupation(t *testing.T, expected, found *RapleafOccupation) {
  if expected == found {
    return
  }
  if expected == nil || found == nil {
    t.Errorf("Expected \n <<<%s>>>\n but found \n<<<%s>>>\n", expected, found)
    return
  }
  if expected.Company != found.Company {
    t.Errorf("Expected company %s but found %s in occupation", expected.Company, found.Company)
  }
  if expected.JobTitle != found.JobTitle {
    t.Errorf("Expected job title %s but found %s in occupation", expected.JobTitle, found.JobTitle)
  }
} 

func testSameMemberSite(t *testing.T, expected, found *RapleafMemberSite) {
  if expected == found {
    return
  }
  if expected == nil || found == nil {
    t.Errorf("Expected \n <<<%s>>>\n but found \n<<<%s>>>\n", expected, found)
    return
  }
  if expected.Site != found.Site {
    t.Errorf("Expected site %s but found %s in membership", expected.Site, found.Site)
  }
  if expected.ProfileUrl != found.ProfileUrl {
    t.Errorf("Expected profile url %s but found %s in membership", expected.ProfileUrl, found.ProfileUrl)
  }
  if expected.ImageUrl != found.ImageUrl {
    t.Errorf("Expected image url %s but found %s in membership", expected.ImageUrl, found.ImageUrl)
  }
  if expected.NumFriends != found.NumFriends {
    t.Errorf("Expected num friends %d but found %d in membership", expected.NumFriends, found.NumFriends)
  }
  if expected.NumFollowers != found.NumFollowers {
    t.Errorf("Expected num followers %d but found %d in membership", expected.NumFollowers, found.NumFollowers)
  }
  if expected.NumFollowed != found.NumFollowed {
    t.Errorf("Expected num followed %d but found %d in membership", expected.NumFollowed, found.NumFollowed)
  }
  if expected.Exists != found.Exists {
    t.Errorf("Expected exists %s but found %s in membership", expected.Exists, found.Exists)
  }
} 

func testSamePerson(t *testing.T, expected, found *RapleafPerson) {
  if expected == found {
    return
  }
  if expected == nil {
    t.Error("Expected \n <<<nil>>>\n but found \n <<<", found, ">>>\n")
    return
  }
  if found == nil {
    t.Error("Expected \n <<<", expected, ">>>\n but found \n <<<nil>>>\n")
    return
  }
  if expected.Id != found.Id {
    t.Errorf("Expected id %s but found %s in person", expected.Id, found.Id)
  }
  if expected.Name != found.Name {
    t.Errorf("Expected name %s but found %s in person", expected.Name, found.Name)
  }
  if expected.Gender != found.Gender {
    t.Errorf("Expected gender %s but found %s in person", expected.Gender, found.Gender)
  }
  if expected.Location != found.Location {
    t.Errorf("Expected location %s but found %s in person", expected.Location, found.Location)
  }
  if expected.NumFriends != found.NumFriends {
    t.Errorf("Expected num friends %d but found %d in person", expected.NumFriends, found.NumFriends)
  }
  if expected.Age != found.Age {
    t.Errorf("Expected age %d but found %d in person", expected.Age, found.Age)
  }
  if expected.EmailAddress != found.EmailAddress {
    t.Errorf("Expected email address %s but found %s in person", expected.EmailAddress, found.EmailAddress)
  }
  if expected.EarliestKnownActivity != found.EarliestKnownActivity {
    if expected.EarliestKnownActivity == nil || found.EarliestKnownActivity == nil || expected.EarliestKnownActivity.Seconds() != found.EarliestKnownActivity.Seconds() {
      t.Errorf("Expected EarliestKnownActivity %v but found %v in person", expected.EarliestKnownActivity, found.EarliestKnownActivity)
    }
  }
  if expected.LatestKnownActivity != found.LatestKnownActivity {
    if expected.LatestKnownActivity == nil || found.LatestKnownActivity == nil || expected.LatestKnownActivity.Seconds() != found.LatestKnownActivity.Seconds() {
      t.Errorf("Expected LatestKnownActivity %v but found %v in person", expected.LatestKnownActivity, found.LatestKnownActivity)
    }
  }
  if len(expected.Occupations) != len(found.Occupations) {
    t.Errorf("Expected %d occupations but found %d occupations in person", len(expected.Occupations), len(found.Occupations))
  } else {
    for i, occupation := range expected.Occupations {
      testSameOccupation(t, occupation, found.Occupations[i])
    }
  }
  if len(expected.Memberships) != len(found.Memberships) {
    t.Errorf("Expected %d memberships but found %d memberships in person", len(expected.Memberships), len(found.Memberships))
  } else {
    for i, membership := range expected.Memberships {
      testSameMemberSite(t, membership, found.Memberships[i])
    }
  }
}

func TestPersonEmptyXml(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  u, err := RapleafPersonFromString(USER_EMPTY_XML)
  closeServerTestFiles(l)
  if err != nil {
    t.Error("Unable to parse empty rapleaf user xml: ", err.String())
    return
  }
  if u == nil {
    t.Error("Unable to parse empty rapleaf user xml")
    return
  }
  testSamePerson(t, USER_EMPTY_PERSON, u)
}

func TestPersonWithProfileXml(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  u, err := RapleafPersonFromString(USER_WITH_PROFILE_XML)
  closeServerTestFiles(l)
  if err != nil {
    t.Error("Unable to parse rapleaf user xml: ", err.String())
    return
  }
  if u == nil {
    t.Error("Unable to parse rapleaf user xml")
    return
  }
  testSamePerson(t, USER_WITH_PROFILE_PERSON, u)
}

func TestPersonXmlByEmail(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  code, text := PersonXmlByEmail(API_KEY, "empty.profile@gmail.com")
  if code != 200 {
    t.Error("Expected status code 200 but received ", code, " with message: ", text)
    closeServerTestFiles(l)
    return
  }
  text = strings.TrimSpace(text)
  if text != USER_EMPTY_XML {
    t.Error("Expected XML:\n ", len(USER_EMPTY_XML), " vs. ", len(text), "\n<<<", USER_EMPTY_XML, ">>>\n but found \n <<<", text, ">>>")
    closeServerTestFiles(l)
    return
  }
  code, text = PersonXmlByEmail(API_KEY, "john.q.public@gmail.com")
  closeServerTestFiles(l)
  if code != 200 {
    t.Error("Expected status code 200 but received ", code, " with message: ", text)
    return
  }
  text = strings.TrimSpace(text)
  if text != USER_WITH_PROFILE_XML {
    t.Error("Expected XML:\n ", len(USER_WITH_PROFILE_XML), " vs. ", len(text), "\n<<<", USER_WITH_PROFILE_XML, ">>>\n but found \n <<<", text, ">>>")
    return
  }
}

func TestPersonXmlByRapleafId(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  code, text := PersonXmlByRapleafId(API_KEY, "b34282025d7e2c5db6786a8daaab48c7")
  if code != 200 {
    t.Error("Expected status code 200 but received ", code, " with message: ", text)
    closeServerTestFiles(l)
    return
  }
  text = strings.TrimSpace(text)
  if text != USER_EMPTY_XML {
    t.Error("Expected XML:\n ", len(USER_EMPTY_XML), " vs. ", len(text), "\n<<<", USER_EMPTY_XML, ">>>\n but found \n <<<", text, ">>>")
    closeServerTestFiles(l)
    return
  }
  code, text = PersonXmlByRapleafId(API_KEY, "97fc425100000000")
  closeServerTestFiles(l)
  if code != 200 {
    t.Error("Expected status code 200 but received ", code, " with message: ", text)
    return
  }
  text = strings.TrimSpace(text)
  if text != USER_WITH_PROFILE_XML {
    t.Error("Expected XML:\n ", len(USER_WITH_PROFILE_XML), " vs. ", len(text), "\n<<<", USER_WITH_PROFILE_XML, ">>>\n but found \n <<<", text, ">>>")
    return
  }
}

func TestPersonXmlBySite(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  code, text := PersonXmlBySite(API_KEY, "linkedin", "johnqpublic")
  closeServerTestFiles(l)
  if code != 200 {
    t.Error("Expected status code 200 but received ", code, " with message: ", text)
    return
  }
  text = strings.TrimSpace(text)
  if text != USER_WITH_PROFILE_XML {
    t.Error("Expected XML:\n ", len(USER_WITH_PROFILE_XML), " vs. ", len(text), "\n<<<", USER_WITH_PROFILE_XML, ">>>\n but found \n <<<", text, ">>>")
    return
  }
}

func TestPersonByEmail(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  expected1 := *USER_EMPTY_PERSON
  expected1.EmailAddress = "empty.profile@gmail.com"
  expected2 := *USER_WITH_PROFILE_PERSON
  expected2.EmailAddress = "john.q.public@gmail.com"
  testSamePerson(t, &expected1, PersonByEmail(API_KEY, "empty.profile@gmail.com"))
  testSamePerson(t, &expected2, PersonByEmail(API_KEY, "john.q.public@gmail.com"))
  closeServerTestFiles(l)
}

func TestPersonByRapleafId(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  testSamePerson(t, USER_EMPTY_PERSON, PersonByRapleafId(API_KEY, "b34282025d7e2c5db6786a8daaab48c7"))
  testSamePerson(t, USER_WITH_PROFILE_PERSON, PersonByRapleafId(API_KEY, "97fc425100000000"))
  closeServerTestFiles(l)
}

func TestPersonBySite(t *testing.T) {
  l, err := serveTestFiles(t)
  if err != nil {
    return
  }
  testSamePerson(t, USER_WITH_PROFILE_PERSON, PersonBySite(API_KEY, "linkedin", "johnqpublic"))
  closeServerTestFiles(l)
}
