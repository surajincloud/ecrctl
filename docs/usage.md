# Usage

* Login to ECR

```
ecrctl login
```

* Login to ECR in different region

```
ecrctl login --region <region>
```

Note: You can also set region in environment variable `AWS_REGION`

* list repositories

```
ecrctl get repositories
```

* list repositories along with their tags

```
ecrctl get repositories --show-tags
```

* filter repositories based on tag

```
ecrctl get repositories --tag key=value
```

you can also use short flag `-t`

```
ecrctl get repositories -t key=value
```

* List images from a given repository

```
ecrctl get images --repo <repo-name>
```
