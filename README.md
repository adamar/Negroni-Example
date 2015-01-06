Negroni Example
===============

This is an example Negroni web app with basic auth


Note
----------
- Please note that this example at the current time lacks a few important security features (ie. password hashing) 
and as such should not be used in production without a lot of changes.

- This example webapp was originally built with a MySQL DB, but has now been modified to use Postgres. There
is now a 'mysql' branch that still uses the go-mysql lib, but development on this branch has slowed to little more
than a trickle



Features
----------
* Signup/Login
* Pages rendered from templates
* Very Simple API call via Javascript on homepage


Requirements
-----------

* Negroni
* Negroni-sessions
* Render
* pq


Try Me
-----------

Deploy to Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/adamar/negroni-example)



To Do
-----------

* Hash passwords
* More extensive api examples


