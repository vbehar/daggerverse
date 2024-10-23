# Git Info Dagger Module

This is a [Dagger](https://dagger.io/) module to extract information from a Git repository.

Give it a (local) Git repository path, and it will return the following information:
- the current branch name
- the current tag name (if any)
- the current commit hash
- the current commit author name
- the current commit time
- the current commit message
- the current "version" (from [git-describe](https://git-scm.com/docs/git-describe))
- the URL of the remote repository (if any)
- the name of the repository (the last part of the URL)

You can retrieve these information:
- in JSON format - as a string output
- in a Dagger File object, in JSON format
- in a Dagger Directory object, 1 file per information, in raw format
- as environment variables, injected in a provided container

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/git-info>.

## Usage

Display the information in JSON format:

```bash
$ dagger call -m github.com/vbehar/daggerverse/git-info \
	--git-directory=. \
	json
{
  "Ref": "HEAD",
  "Branch": "main",
  "Tag": "v1.2.3",
  "CommitHash": "5735dead51a44ea8d954508d8e36b7facf2f49dd",
  "CommitUser": "Vincent Behar",
  "CommitTime": "2024-10-22T11:16:27+02:00",
  "CommitMessage": "feat: git-info",
  "Version": "v1.2.3",
  "URL": "https://github.com/vbehar/daggerverse",
  "Name": "daggerverse"
}
```
