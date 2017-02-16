package testdata

var IssueResponse = []byte(`
[
    {
        "assignee": null,
        "assignees": [],
        "body": "The UDIV and SDIV instructions are optional on ARM, so arm gcc generates __armeabi_udiv(a, b) for a/b by default, but it also emits \"udiv a, b\" while -march=armv7ve is specified.\r\n\r\nGolang should also allow user to choose hardware or software division. Maybe by adding a GOARMHDIV environment variable?\r\n\r\n\r\n\r\n\r\n",
        "closed_at": null,
        "comments": 9,
        "comments_url": "https://api.github.com/repos/golang/go/issues/19118/comments",
        "created_at": "2017-02-16T02:45:11Z",
        "events_url": "https://api.github.com/repos/golang/go/issues/19118/events",
        "html_url": "https://github.com/golang/go/issues/19118",
        "id": 207997623,
        "labels": [],
        "labels_url": "https://api.github.com/repos/golang/go/issues/19118/labels{/name}",
        "locked": false,
        "milestone": null,
        "number": 19118,
        "reactions": {
            "+1": 0,
            "-1": 0,
            "confused": 0,
            "heart": 0,
            "hooray": 0,
            "laugh": 0,
            "total_count": 0,
            "url": "https://api.github.com/repos/golang/go/issues/19118/reactions"
        },
        "repository_url": "https://api.github.com/repos/golang/go",
        "state": "open",
        "title": "Use SDIV and UDIV for ARM",
        "updated_at": "2017-02-16T04:41:43Z",
        "url": "https://api.github.com/repos/golang/go/issues/19118",
        "user": {
            "avatar_url": "https://avatars.githubusercontent.com/u/24586233?v=3",
            "events_url": "https://api.github.com/users/benshi001/events{/privacy}",
            "followers_url": "https://api.github.com/users/benshi001/followers",
            "following_url": "https://api.github.com/users/benshi001/following{/other_user}",
            "gists_url": "https://api.github.com/users/benshi001/gists{/gist_id}",
            "gravatar_id": "",
            "html_url": "https://github.com/benshi001",
            "id": 24586233,
            "login": "benshi001",
            "organizations_url": "https://api.github.com/users/benshi001/orgs",
            "received_events_url": "https://api.github.com/users/benshi001/received_events",
            "repos_url": "https://api.github.com/users/benshi001/repos",
            "site_admin": false,
            "starred_url": "https://api.github.com/users/benshi001/starred{/owner}{/repo}",
            "subscriptions_url": "https://api.github.com/users/benshi001/subscriptions",
            "type": "User",
            "url": "https://api.github.com/users/benshi001"
        }
    },
    {
        "assignee": {
            "avatar_url": "https://avatars.githubusercontent.com/u/2688315?v=3",
            "events_url": "https://api.github.com/users/aclements/events{/privacy}",
            "followers_url": "https://api.github.com/users/aclements/followers",
            "following_url": "https://api.github.com/users/aclements/following{/other_user}",
            "gists_url": "https://api.github.com/users/aclements/gists{/gist_id}",
            "gravatar_id": "",
            "html_url": "https://github.com/aclements",
            "id": 2688315,
            "login": "aclements",
            "organizations_url": "https://api.github.com/users/aclements/orgs",
            "received_events_url": "https://api.github.com/users/aclements/received_events",
            "repos_url": "https://api.github.com/users/aclements/repos",
            "site_admin": false,
            "starred_url": "https://api.github.com/users/aclements/starred{/owner}{/repo}",
            "subscriptions_url": "https://api.github.com/users/aclements/subscriptions",
            "type": "User",
            "url": "https://api.github.com/users/aclements"
        },
        "assignees": [
            {
                "avatar_url": "https://avatars.githubusercontent.com/u/2688315?v=3",
                "events_url": "https://api.github.com/users/aclements/events{/privacy}",
                "followers_url": "https://api.github.com/users/aclements/followers",
                "following_url": "https://api.github.com/users/aclements/following{/other_user}",
                "gists_url": "https://api.github.com/users/aclements/gists{/gist_id}",
                "gravatar_id": "",
                "html_url": "https://github.com/aclements",
                "id": 2688315,
                "login": "aclements",
                "organizations_url": "https://api.github.com/users/aclements/orgs",
                "received_events_url": "https://api.github.com/users/aclements/received_events",
                "repos_url": "https://api.github.com/users/aclements/repos",
                "site_admin": false,
                "starred_url": "https://api.github.com/users/aclements/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/aclements/subscriptions",
                "type": "User",
                "url": "https://api.github.com/users/aclements"
            }
        ],
        "body": "Currently the GC doesn't always wake up idle Ps, and hence may not take full advantage of idle marking during the concurrent mark phase. This can happen during mark 2 because mark 1 completion preempts all workers; if the Ps running those workers have nothing else to do they will simply park, and there's no mechanism to wake them up after we allow mark workers to start again. It's possible this can happen during mark 1 as well, though it may be since we allow mark workers to run before starting the world for mark 1 that all of the Ps start running.\n\n/cc @RLH \n",
        "closed_at": null,
        "comments": 8,
        "comments_url": "https://api.github.com/repos/golang/go/issues/14179/comments",
        "created_at": "2016-02-01T19:38:31Z",
        "events_url": "https://api.github.com/repos/golang/go/issues/14179/events",
        "html_url": "https://github.com/golang/go/issues/14179",
        "id": 130460512,
        "labels": [
            {
                "color": "ededed",
                "default": false,
                "id": 373399998,
                "name": "NeedsFix",
                "url": "https://api.github.com/repos/golang/go/labels/NeedsFix"
            }
        ],
        "labels_url": "https://api.github.com/repos/golang/go/issues/14179/labels{/name}",
        "locked": false,
        "milestone": {
            "closed_at": null,
            "closed_issues": 18,
            "created_at": "2016-09-29T18:06:31Z",
            "creator": {
                "avatar_url": "https://avatars.githubusercontent.com/u/115761?v=3",
                "events_url": "https://api.github.com/users/quentinmit/events{/privacy}",
                "followers_url": "https://api.github.com/users/quentinmit/followers",
                "following_url": "https://api.github.com/users/quentinmit/following{/other_user}",
                "gists_url": "https://api.github.com/users/quentinmit/gists{/gist_id}",
                "gravatar_id": "",
                "html_url": "https://github.com/quentinmit",
                "id": 115761,
                "login": "quentinmit",
                "organizations_url": "https://api.github.com/users/quentinmit/orgs",
                "received_events_url": "https://api.github.com/users/quentinmit/received_events",
                "repos_url": "https://api.github.com/users/quentinmit/repos",
                "site_admin": false,
                "starred_url": "https://api.github.com/users/quentinmit/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/quentinmit/subscriptions",
                "type": "User",
                "url": "https://api.github.com/users/quentinmit"
            },
            "description": "",
            "due_on": "2017-04-30T07:00:00Z",
            "html_url": "https://github.com/golang/go/milestone/47",
            "id": 2038341,
            "labels_url": "https://api.github.com/repos/golang/go/milestones/47/labels",
            "number": 47,
            "open_issues": 129,
            "state": "open",
            "title": "Go1.9Early",
            "updated_at": "2017-02-15T22:43:37Z",
            "url": "https://api.github.com/repos/golang/go/milestones/47"
        },
        "number": 14179,
        "reactions": {
            "+1": 0,
            "-1": 0,
            "confused": 0,
            "heart": 0,
            "hooray": 0,
            "laugh": 0,
            "total_count": 0,
            "url": "https://api.github.com/repos/golang/go/issues/14179/reactions"
        },
        "repository_url": "https://api.github.com/repos/golang/go",
        "state": "open",
        "title": "runtime: GC should wake up idle Ps",
        "updated_at": "2017-02-16T03:27:53Z",
        "url": "https://api.github.com/repos/golang/go/issues/14179",
        "user": {
            "avatar_url": "https://avatars.githubusercontent.com/u/2688315?v=3",
            "events_url": "https://api.github.com/users/aclements/events{/privacy}",
            "followers_url": "https://api.github.com/users/aclements/followers",
            "following_url": "https://api.github.com/users/aclements/following{/other_user}",
            "gists_url": "https://api.github.com/users/aclements/gists{/gist_id}",
            "gravatar_id": "",
            "html_url": "https://github.com/aclements",
            "id": 2688315,
            "login": "aclements",
            "organizations_url": "https://api.github.com/users/aclements/orgs",
            "received_events_url": "https://api.github.com/users/aclements/received_events",
            "repos_url": "https://api.github.com/users/aclements/repos",
            "site_admin": false,
            "starred_url": "https://api.github.com/users/aclements/starred{/owner}{/repo}",
            "subscriptions_url": "https://api.github.com/users/aclements/subscriptions",
            "type": "User",
            "url": "https://api.github.com/users/aclements"
        }
    },
    {
        "assignee": null,
        "assignees": [],
        "body": "Please answer these questions before submitting your issue. Thanks!\r\n\r\n### What version of Go are you using (go version)?\r\ngolang.org/x/arch/arm\r\n\r\n### What operating system and processor architecture are you using (go env)?\r\nUbuntu 16.04.1 LTS and ARM64\r\n\r\n### What did you do?\r\nDisassemble following two instructions\r\n937facb1\r\nd6530d61\r\n\r\n### What did you expect to see?\r\nSTREXD.LT [R12], R4, R3, R7\r\nLDRD.VS [SP, -R6], R6, R5\r\n\r\n### What did you see instead?\r\nSTREXD.LT [R12], R3, R3, R7\r\nLDRD.VS [SP, -R6], R5, R5\r\n\r\n",
        "closed_at": null,
        "comments": 4,
        "comments_url": "https://api.github.com/repos/golang/go/issues/19100/comments",
        "created_at": "2017-02-15T07:42:55Z",
        "events_url": "https://api.github.com/repos/golang/go/issues/19100/events",
        "html_url": "https://github.com/golang/go/issues/19100",
        "id": 207729562,
        "labels": [],
        "labels_url": "https://api.github.com/repos/golang/go/issues/19100/labels{/name}",
        "locked": false,
        "milestone": null,
        "number": 19100,
        "reactions": {
            "+1": 0,
            "-1": 0,
            "confused": 0,
            "heart": 0,
            "hooray": 0,
            "laugh": 0,
            "total_count": 0,
            "url": "https://api.github.com/repos/golang/go/issues/19100/reactions"
        },
        "repository_url": "https://api.github.com/repos/golang/go",
        "state": "open",
        "title": "x/arch/arm/armasm: second source register is calculated incorrectly",
        "updated_at": "2017-02-16T03:08:37Z",
        "url": "https://api.github.com/repos/golang/go/issues/19100",
        "user": {
            "avatar_url": "https://avatars.githubusercontent.com/u/23331893?v=3",
            "events_url": "https://api.github.com/users/williamweixiao/events{/privacy}",
            "followers_url": "https://api.github.com/users/williamweixiao/followers",
            "following_url": "https://api.github.com/users/williamweixiao/following{/other_user}",
            "gists_url": "https://api.github.com/users/williamweixiao/gists{/gist_id}",
            "gravatar_id": "",
            "html_url": "https://github.com/williamweixiao",
            "id": 23331893,
            "login": "williamweixiao",
            "organizations_url": "https://api.github.com/users/williamweixiao/orgs",
            "received_events_url": "https://api.github.com/users/williamweixiao/received_events",
            "repos_url": "https://api.github.com/users/williamweixiao/repos",
            "site_admin": false,
            "starred_url": "https://api.github.com/users/williamweixiao/starred{/owner}{/repo}",
            "subscriptions_url": "https://api.github.com/users/williamweixiao/subscriptions",
            "type": "User",
            "url": "https://api.github.com/users/williamweixiao"
        }
    }
]`)
