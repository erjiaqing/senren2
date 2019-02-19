# ProblemCI 2

Problem CI rewritten in Golang

- Better ACL
  - Only owner can access problem by UID, others should access problem by access key
  - Owner can create, update, review and revoke access key at any time

- Git server is no longer public, it will be shipped with docker compose (gitea)
  - You can access some gitea api by our proxy, with strict permission
  - For example, you should access it by /giteaapi/?api={gitea-api-suffix}&method={gitea-api-method}&key={your access key}&time={time}&sign={sha1(private-key + ".." + gitea-api-suffix + ".." + gitea-api-method + ".." + unix time stamp in hex)}
  - Time error should be less than 30 seconds

- Now, git server is only used as a storage, you should invoke build task by your self