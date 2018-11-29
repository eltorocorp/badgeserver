# BadgesServer
BadgeServer is a tiny server that allows private repositories to use shields.io for status badges.
![Build Status](http://badges.awsp.eltoro.com?project=badgeserver&item=build)

## How it works
A BadgeServer instance sits as an intermediary between a repository (i.e. github), a CI Server (i.e. Jenkins), and [shields.io](shelds.io), allowing them to communicate badge information without knowing about one another.

When the CI server completes a build, it sends a POST request to the BadgeServer with some badge information.

A repository contributor places a link on the github readme page where a badge is desired. The link points to the badges server, and includes query parameters which denote which badge information to retrieve.

When github renders the link into the page, a GET request is sent to the badge server. The badge server looks up the badge information, and returns the requested badge image.

## Example
### Step 1
Let's assume we have some project called `myproject`.

The following markup is on the `myproject` project's README.md:
```
[![build status](http://badgeserver.mydomain.com/?project=myproject&item=build_status)](http://www.somelink.com)
```
### Step 2

A build completes successfully and sends a POST request to the BadgeServer. The post request includes a query string with parameters that describe the kind of badge to generate.

 - **project**: The name of the project to create the badge for. i.e. "myproject"
 - **item**: The name of the badge to create. i.e. "build_status"
 - **value**: A value for the badge. i.e. "passing"
 - **color**: A color that the badge should be. i.e. "green"

```
POST / HTTP/1.1
Host: http://badgeserver.mydomain.com/
Content-Type: application/x-www-form-urlencoded
Content-Length: 61

project=myproject&item=build_status&value=passing&color=green
```
How the POST request is accomplished is up to the implementor, but it makes sense to just use a Python or Shell script of some sort.

For convenience, several helper scripts have been included in the `helpers` directory.

### Step 3

When github loads the markdown, which was created in step one, it sends a GET request to the BadgeServer requesting the badge resource. The BadgeServer then translates that request to a request sent to [shields.io](https://shields.io/). This is where the actual badge image is generated. The response from shields.io is intercepted by the badges server and modified such that it will work with github repositories (github has an image caching layer that needs to be bypassed). The modified response is then forwarded back to github and github renders the requested image into the page.

## Sure, but why?

A number of badge systems exist that integrate with build servers and github, but, in this engineer's opinion, it's all just too complicated and restrictive.

Things get dicy when you have a badge that a service doesn't exist to represent. 

Things get even more complicated when you want to integrate some of these systems with private repositories, particularly when you don't necessarily trust the other parties or want them to know anything about your repository.

Shields.io is nice, but has caching limitations when it comes to ad-hoc badge generation. BadgeServer leverages shields.io to generate badges, but fixes the caching behavior. (Honestly, the caching limitations could be overcome pretty easily, but who wants to sit around waiting for a pull request to get merged? Not me.)

Lastly, it is worth noting that the BadgeServer will work with any system requesting a badge image (not just github and jenkins).