## clnotifications

The clnotifications tool cleans up thousands of Windows Push Notification Platform
 (WPN) and Windows Notification Facility (WNF) values from `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Notifications` registry key. To prevent high CPU usage and slow user logons.

By default the tool skip first 500 values ans clean up the rest. It puts the logs into the same directory as the execution. 

```term
 & .\clnotifications.exe -?
clnotifications v0.1.1
flag provided but not defined: -?
Usage of C:\Users\kurlo\Documents\golang\clnotifications\clnotifications.exe:
  -cleanup
        command to start cleaning up (default true)
  -count-values-in-chunks int
        number of values to delete in one chunk (default 100)
  -count-values-to-read int
        number of values to read in one iteration (default 1000)
  -count-values-to-skip-key int
        number of values to skip deletion (default 500)
```

---

* There is a correlation between amount of values at `Notifications` the registry key and long logon time.
  * Microsoft has acknowledged that the registry key `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Notifications` can get bloated with a large number of entries, which can cause issues such as slow logons. This happens due to leaked WNF registrations.
* The `explorer.exe` and `frxsvc.exe` reads thousands of WPN and WNF values from this registry key during user logon. This can lead to high CPU usage and slow logons. See: [~30 second delay on the "Please wait for the FSLogix Apps Services" portion of the login](https://learn.microsoft.com/en-us/archive/msdn-technet-forums/b3a2e9d9-b073-44d9-aea4-9792e0095b4c)
* To resolve this issue, Microsoft provides a fix for Windows Server 2012 R2 that prevents WNF and WNF registrations from being leaked after its installation. The fix also includes a tool called `wnfcleanup` that removes stale WNF registrations created before the installation of the leak fix.
* There is no official input from Microsoft on how to fix this issue for Windows Server 2019.

#### Reference
* [Registry bloat causes slow logons or insufficient system resources error 0x800705AA in Windows 8.1](https://support.microsoft.com/en-us/topic/registry-bloat-causes-slow-logons-or-insufficient-system-resources-error-0x800705aa-in-windows-8-1-82a985fb-df27-abda-440b-f3f81a2c949d)
* https://msendpointmgr.com/2021/06/17/fslogix-slow-sign-in-fix/
* https://answers.microsoft.com/en-us/windows/forum/all/how-to-disable-everything-related-to-windows-push/5b9522ad-cebd-47b6-b7c5-1620da167f45
* https://superuser.com/questions/1246950/can-windows-push-notification-services-wns-be-disabled
* https://www.tenable.com/audits/items/CIS_MS_Windows_11_Enterprise_Level_2_Next_Generation_Windows_Security_v1.0.0.audit:01105485ef3edda773199f06cab046ac
