<?php
/*
 * $Header:   $
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
  
  // specifies member sites with all info for an existing profile
  // not all fields will be populated (especially number fields, defaulting to 0)
  class RapleafMemberSite {
    public $profile_url;
    public $image_url;
    public $num_friends;
    public $num_followers;
    public $num_followed;
    
    public function __construct($membership) {
      $this->profile_url = (string)$membership["profile_url"] or '';
      $this->num_friends = (int)((string)$membership["num_friends"]);
      $this->num_followers = (int)((string)$membership["num_followers"]);
      $this->num_followed = (int)((string)$membership["num_followed"]);
      $this->image_url = (string)$membership["image_url"] or '';
    }
  }
  
  // specifies user information for a general user
  class RapleafUser {
    public $rapleaf_id;
    public $name;
    public $gender;
    public $location;
    public $num_friends;
    public $occupation_title;
    public $occupation_company;
    public $profiles;
    public $xml;
    
    public function __construct($xml_text) {
      $this->rapleaf_id = '';
      $this->name = '';
      $this->gender = '';
      $this->location = '';
      $this->occupation_title = '';
      $this->profiles = array();
      $this->num_friends = 0;
      $this->xml = $xml_text;
      
      if($xml_text) {
        $xml = new SimpleXMLElement($xml_text);
        $this->rapleaf_id = (string)$xml["id"];
        $this->name = (string)$xml->basics->name or '';
        $this->gender = (string)$xml->basics->gender or '';
        $this->location = (string)$xml->basics->location or '';
        $this->num_friends = (int)((string)$xml->basics->num_friends or 0);
        $this->occupation_title = (string)$xml->basics->occupations->occupation[0]['job_title'][0] or '';
        $this->occupation_company = (string)$xml->basics->occupations->occupation[0]['company'][0] or '';
        foreach($xml->memberships->primary->membership as $membership) {
          if($membership["exists"] == "true") {
            $this->profiles[(string)$membership["site"]] = new RapleafMemberSite($membership);
          }
        }
        foreach($xml->memberships->supplemental->membership as $membership) {
          if($membership["exists"] == "true") {
            $this->profiles[(string)$membership["site"]] = new RapleafMemberSite($membership);
          }
        }
      }
    }
  }
  
  // class for accessing rapleaf data
  class Rapleaf {
    private $apikey;
    private $db;
    private $sites;
    
    public function __construct($apikey="stuff") {
      $this->apikey = $apikey;
      $this->sites = array("bebo" => "Bebo", 
                           "facebook" => "Facebook",
                           "flickr" => "Flickr",
                           "friendster" => "Friendster",
                           "hi5" => "hi5",
                           "linkedin" => "LinkedIn",
                           "myspace" => "MySpace",
                           "plaxo" => "Plaxo",
                           "rapleaf" => "Rapleaf",
                           "twitter" => "Twitter",
      );
    }

    private function dohttprequest($method, $url, $headers, $data) {
      $ch = curl_init();
      curl_setopt($ch, CURLOPT_URL, $url);
      if($method == "POST") {
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
      }
      if($headers) {
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
      }
      curl_setopt($ch, CURLOPT_HEADER, true);
      curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
      $errno = curl_errno($ch);
      $output = curl_exec($ch);
      curl_close($ch);
      return array($output, $errno);
    }
    
    public function byEmail($email) {
      if(!$email)
        return null;
      $rapleaf_user = null;
      $url = "http://api.rapleaf.com/v3/person/email/" . urlencode($email);
      $res = $this->dohttprequest("GET", $url, array('Authorization: ' . $this->apikey), false);
      $response = $res[0];
      $errno = $res[1];
      if(!$errno && $response) {
        $rapleaf_user = new RapleafUser($response);
      }
      return $rapleaf_user;
    }
    
    public function bySiteAndProfileId($site, $profile_id) {
      if(!$site || !$profile_id)
        return null;
      $rapleaf_user = null;
      $url = "http://api.rapleaf.com/v3/person/web/" . urlencode($site) . "/" . urlencode($profile_id);
      $res = $this->dohttprequest("GET", $url, array('Authorization: ' . $this->apikey), false);
      $response = $res[0];
      $errno = $res[1];
      if(!$errno && $response) {
        $rapleaf_user = new RapleafUser($response);
      }
      return $rapleaf_user;
    }
    
    public function retrieveUser($email = '', $site = '', $profile_id = '') {
      $email = (string)$email or "";
      $site = (string)$site or "";
      $profile_id = (string)$profile_id or "";
      $xml = null;
      $was_cached = false;
      
      if($email) {
        $rapleaf_user = $this->byEmail($email);
      } else if($site && $profile_id) {
        $rapleaf_user = $this->bySiteAndProfileId($site, $profile_id);
      }
      return $rapleaf_user;
    }
    
    public function retrieveUserIds($emailOrUserId) {
      $url = "http://api.rapleaf.com/v2/graph/" . urlencode($emailOrUserId) . "?n=1";
      $res = $this->dohttprequest("GET", $url, array('Authorization: ' . $this->apikey), false);
      $response = $res[0];
      $errno = $res[1];
      if(!$errno && $response) {
        return explode("\n", $response);
      }
      return array();
    }
    
    public function retrieveEmailAddresses($emailOrUserId) {
      $url = "http://api.rapleaf.com/v2/graph/" . urlencode($emailOrUserId) . "?n=2";
      $res = $this->dohttprequest("GET", $url, array('Authorization: ' . $this->apikey), false);
      $response = $res[0];
      $errno = $res[1];
      if(!$errno && $response) {
        return explode(",", $response);
      }
      return array();
    }
  }
?>