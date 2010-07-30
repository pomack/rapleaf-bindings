<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html>
<!--
  # $Header:   $
  # Copyright 2010 Aalok Shah (aalok@shah.ws)
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
-->
<head>
  <title>Find Info from Rapleaf</title>
</head>
<body>
  <h1>Find Info from Rapleaf</h1>
  <form name="search" method="GET" action="rapleaf_test.php">
    <div><label for="api_key">API Key</label><input name="api_key" value="<?php echo $_GET["api_key"]; ?>"/></div>
    <div>Search by</div>
    <div>
      <label for="email">Email</label><input name="email" value="<?php echo $_GET["email"]; ?>"/>
    </div>
    <div>
      <label for="site">Site</label>
      <select name="site">
        <option value="bebo"<?php if($_GET["site"] == "bebo") { ?> selected="selected"<?php } ?>>Bebo</option>
        <option value="facebook"<?php if($_GET["site"] == "facebook") { ?> selected="selected"<?php } ?>>Facebook</option>
        <option value="flickr"<?php if($_GET["site"] == "flickr") { ?> selected="selected"<?php } ?>>Flickr</option>
        <option value="friendster"<?php if($_GET["site"] == "friendster") { ?> selected="selected"<?php } ?>>Friendster</option>
        <option value="hi5"<?php if($_GET["site"] == "hi5") { ?> selected="selected"<?php } ?>>Hi5</option>
        <option value="linkedin"<?php if($_GET["site"] == "linkedin") { ?> selected="selected"<?php } ?>>LinkedIn</option>
        <option value="myspace"<?php if($_GET["site"] == "myspace") { ?> selected="selected"<?php } ?>>MySpace</option>
        <option value="plaxo"<?php if($_GET["site"] == "plaxo") { ?> selected="selected"<?php } ?>>Plaxo</option>
        <option value="rapleaf"<?php if($_GET["site"] == "rapleaf") { ?> selected="selected"<?php } ?>>Rapleaf</option>
        <option value="twitter"<?php if($_GET["site"] == "twitter") { ?> selected="selected"<?php } ?>>Twitter</option>
      </select>
      <label for="profile_id">Profile Id</label><input name="profile_id" value="<?php echo $_GET["profile_id"]; ?>"/>
    </div>
    <input type="submit" value="submit"/>
  </form>
<?php
  include("rapleaf.php"); 
  $rapleaf = new RapleafAPI($_GET["api_key"]);
  $rapleaf_user = $rapleaf->retrieveUser($_GET["email"], $_GET["site"], $_GET["profile_id"]);
?>
  <div>Rapleaf User: <?php print_r($rapleaf_user); ?></div>
  <div>Rapleaf ID: <?php if($rapleaf_user) {echo $rapleaf_user->rapleaf_id; } ?></div>
  <div>Rapleaf XML: <?php if($rapleaf_user) {echo $rapleaf_user->xml; } ?></div>
</body>
</html>
