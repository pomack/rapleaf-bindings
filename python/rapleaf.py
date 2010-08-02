#!/usr/bin/env python

# Copyright 2010 Aalok Shah
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#      http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Needed to avoid ambiguity in imports
from __future__ import absolute_import

import xml.etree.ElementTree as ET
from urllib import quote
import urllib2
import datetime

__site_profile_id_to_profile_url_template = {
  'bebo'       : 'http://www.bebo.com/%s',
  'facebook'   : 'http://www.facebook.com/%s',
  'flickr'     : 'http://www.flickr.com/%s',
  'friendster' : 'http://profiles.friendster.com/%s',
  'hi5'        : 'http://www.hi5.com/friend/%s',
  'linkedin'   : 'http://www.linkedin.com/in/%s',
  'myspace'    : 'http://profile.myspace.com/%s',
  'plaxo'      : 'http://www.plaxocom/%s',
  'rapleaf'    : 'http://api.rapleaf.com/v3/person/web/rapleaf/%s',
  'twitter'    : 'http://www.twitter.com/%s',
}

def site_profile_id_to_profile_url(site, profile_id):
  global __site_profile_id_to_profile_url_template
  if site in __site_profile_id_to_profile_url_template:
    return __site_profile_id_to_profile_url_template[site] % quote(unicode(profile_id))
  return None

def _int_or_none(value):
  try:
    return int(value) if value is not None else None
  except:
    return None

def _str_to_date(value):
  if value and isinstance(value, basestring) and '-' in value:
    the_date = [int(x, 10) for x in value.split('-', 3)]
    if len(the_date) == 3:
      return datetime.date(*the_date)
  return value or None

class RapleafMemberSite(object):
  def __init__(self, profile_url=None, image_url=None, num_friends=None, num_followers=None, num_followed=None):
    self.profile_url = profile_url or None
    self.image_url = image_url or None
    self.num_friends = _int_or_none(num_friends)
    self.num_followers = _int_or_none(num_followers)
    self.num_followed = _int_or_none(num_followed)
  
  def __repr__(self):
    d = {
      'profile_url' : self.profile_url,
      'image_url' : self.image_url,
      'num_friends' : self.num_friends,
      'num_followers' : self.num_followers,
      'num_followed' : self.num_followed,
    }
    return u"RapleafMemberSite(%s)" % (u','.join([u"%s=%r" % (k,v) for k,v in sorted(d.items()) if v is not None]))
  
  def __unicode__(self):
    d = {
      'profile_url' : self.profile_url,
      'image_url' : self.image_url,
      'num_friends' : self.num_friends,
      'num_followers' : self.num_followers,
      'num_followed' : self.num_followed,
    }
    for k,v in d.items():
      if v is None:
        del d[k]
    return u"RapleafMemberSite%r" % d
  
  @classmethod
  def from_xml(cls, elem):
    if elem is None:
      return None
    if elem and isinstance(elem, basestring):
      elem = ET.XML(elem)
    if elem.attrib and elem.attrib.get('exists') == 'true':
      return RapleafMemberSite(profile_url=elem.attrib.get('profile_url'), 
                               image_url=elem.attrib.get('image_url'), 
                               num_friends=elem.attrib.get('num_friends'), 
                               num_followers=elem.attrib.get('num_followers'), 
                               num_followed=elem.attrib.get('num_followed'))
    return None
  

class RapleafOccupation(object):
  def __init__(self, company=None, job_title=None):
    self.company = company
    self.job_title = job_title
  
  def __repr__(self):
    d = {
      'company' : self.company,
      'job_title' : self.job_title,
    }
    return "RapleafOccupation(%s)" % (u','.join([u"%s=%r" % (k,v) for k,v in sorted(d.items()) if v is not None]))

  def __unicode__(self):
    d = {
      'company' : self.company,
      'job_title' : self.job_title,
    }
    for k,v in d.items():
      if v is None:
        del d[k]
    return "RapleafOccupation%r" % d
  
  @classmethod
  def from_xml(cls, elem):
    if elem is None:
      return None
    if elem and isinstance(elem, basestring):
      elem = ET.XML(elem)
    if elem.tag == 'occupation':
      return RapleafOccupation(elem.attrib.get('company'), elem.attrib.get('job_title'))
    if elem.tag == 'occupations':
      return [RapleafOccupation(e.attrib.get('company'), e.attrib.get('job_title')) for e in elem]
    return None
  

class RapleafUser(object):
  def __init__(self, rapleaf_id=None, name=None, gender=None, location=None, num_friends=None, age=None, earliest_known_activity=None, latest_known_activity=None, email_address=None, occupations=[], universities=[], profiles=[]):
    self.rapleaf_id = rapleaf_id or None
    self.name = name or None
    self.gender = gender.lower() if gender else None 
    self.location = location or None
    self.num_friends = _int_or_none(num_friends)
    self.age = _int_or_none(age)
    self.earliest_known_activity = _str_to_date(earliest_known_activity)
    self.latest_known_activity = _str_to_date(latest_known_activity)
    self.email_address = email_address or None
    self.occupations = occupations or []
    self.universities = universities or []
    self.profiles = profiles or []
  
  def __repr__(self):
    d = {
      'rapleaf_id' : self.rapleaf_id,
      'name' : self.name,
      'gender' : self.gender,
      'location' : self.location,
      'num_friends' : self.num_friends,
      'age' : self.age,
      'earliest_known_activity' : self.earliest_known_activity,
      'latest_known_activity' : self.latest_known_activity,
      'email_address' : self.email_address,
      'occupations' : self.occupations,
      'universities' : self.universities,
      'profiles' : self.profiles,
    }
    return "RapleafUser(%s)" % (u','.join([u"%s=%r" % (k,v) for k,v in sorted(d.items()) if v is not None]))
  
  def __unicode__(self):
    d = {
      'rapleaf_id' : self.rapleaf_id,
      'name' : self.name,
      'gender' : self.gender,
      'location' : self.location,
      'num_friends' : self.num_friends,
      'age' : self.age,
      'earliest_known_activity' : self.earliest_known_activity,
      'latest_known_activity' : self.latest_known_activity,
      'email_address' : self.email_address,
      'occupations' : self.occupations,
      'universities' : self.universities,
      'profiles' : self.profiles,
    }
    for k,v in d.items():
      if v is None or v == []:
        del d[k]
    return "RapleafUser%r" % d
  
  @classmethod
  def from_xml(cls, elem, email_address=None):
    if elem is None:
      return None
    if elem and isinstance(elem, basestring):
      elem = ET.XML(elem)
    if elem.attrib:
      root = elem
      profiles = []
      d = {'email_address' : email_address or None}
      rapleaf_id = root.attrib.get('id')
      if rapleaf_id:
        d['rapleaf_id'] = unicode(rapleaf_id)
      for elem in root:
        if elem.tag == 'basics':
          for subelem in elem:
            tag = str(subelem.tag)
            if tag in ('name', 'gender', 'location', 'earliest_known_activity', 'latest_known_activity', 'age', 'num_friends'):
              d[tag] = subelem.text
            elif tag == 'occupations':
              d['occupations'] = [RapleafOccupation.from_xml(occupation_tag) for occupation_tag in subelem]
            elif tag == 'universities':
              d['universities'] = [university_tag.text for university_tag in subelem]
            else:
              pass
              #print 'unknown tag: %s' % tag
            #elif tag == 'reputation':
            #  for reptag in subelem:
            #    d = {}
            #    if reptag.tag in ('score', 'commerce_score', 'percent_positive', 'profile_url'):
            #      d[reptag.tag] = reptag.text
            #    elif reptag.tag == 'badges':
            #      d['badges'] = dict([(badge.attrib.get('type'), badge.text) for badge in reptag if badge.tag == 'badge'])
            #    self.reputation = d
        elif elem.tag == 'memberships':
          for memberships in elem:
            profiles += [RapleafMemberSite.from_xml(membership) for membership in memberships]
        else:
          pass
          #print 'unknown tag: %s' % elem.tag
      d['profiles'] = [profile for profile in profiles if profile]
      return RapleafUser(**d)
    return None
  

class RapleafAPI(object):
  PERSON_URI = 'http://api.rapleaf.com/v3/person/'
  GRAPH_URI = 'http://api.rapleaf.com/v2/graph/'
  
  ERROR_CODES = {
    200 : 'Request processed successfully.',
    202 : 'This person is currently being searched. Check back shortly and we should have data.',
    400 : 'Invalid email address or Rapleaf ID.',
    401 : 'API key was not provided or is invalid.',
    403 : 'Your query limit has been exceeded. Contact developer@rapleaf.com if you would like to increase your limit.',
    404 : 'Returned for lookup by hash or site userid. We do not have this person in our system. If you would like better results, consider supplying the email address.',
    500 : 'There was an unexpected error on our server. This should be very rare and if you see it please contact developer@rapleaf.com.',
  }
  
  def __init__(self, api_key=None):
    self.__api_key = api_key
  
  def set_api_key(self, api_key):
    self.__api_key = api_key
  
  def api_key(self):
    return self.__api_key
  
  def __get_request(self, url):
    if not self.api_key():
      return '', 400
    if not url:
      return '', 400
    headers = {'Authorization' : self.api_key()}
    try:
      req = urllib2.Request(url, None, headers)
      result = urllib2.urlopen(req)
      return result.read(), 200
    except urllib2.HTTPError, err:
      return '', err.code
  
  def person_by_rapleaf_id(self, rapleaf_id):
    return self.person_by_profile('rapleaf', rapleaf_id)
  
  def person_by_email(self, email):
    url = self.__class__.PERSON_URI + 'email/' + quote(unicode(email))
    value, code = self.__get_request(url)
    if code == 200 and value:
      print value
      return code, None, RapleafUser.from_xml(value, email_address=email), code
    return code, self.__class__.ERROR_CODES.get(code, 'Unknown Error Code'), None
  
  def person_by_profile(self, site, profile_id):
    url = self.__class__.PERSON_URI + 'web/' + unicode(site) + '/' + quote(unicode(profile_id))
    value, code = self.__get_request(url)
    if code == 200 and value:
      print value
      return code, None, RapleafUser.from_xml(value)
    return code, self.__class__.ERROR_CODES.get(code, 'Unknown Error Code'), None
  
  def graph_output_rapleaf_ids(self, rapleaf_id_or_email):
    return self._get_graph(rapleaf_id_or_email, 1)
  
  def graph_output_emails(self, rapleaf_id_or_email):
    return self._get_graph(rapleaf_id_or_email, 2)
  
  def _get_graph(self, rapleaf_id_or_email, n=1):
    n = _int_or_none(n) or 1
    url = self.__class__.GRAPH_URI + quote(unicode(rapleaf_id_or_email) + '?' + unicode(n))
    value, code = self.__get_request(url)
    if code == 200 and value:
      print value
      if n == 2:
        return code, None, value.split(',')
      else:
        return code, None, value.split('\n')
    return code, self.__class__.ERROR_CODES.get(code, 'Unknown Error Code'), None
  

