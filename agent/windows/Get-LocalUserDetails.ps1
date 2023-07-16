param (
    [String] $UserNames
)

$userNamesArray = $UserNames -split ','

$localUsersDetail = @()

$users = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True'"
foreach ($userName in $userNamesArray) {
    foreach ($user in $users) {
        if ($UserName -eq "" -or $user.Name -eq $UserName) {
            $localUsersDetail += @{
                'Name'               = $user.Name
                'SID'                = $user.SID
                'FullName'           = $user.FullName
                'Description'        = $user.Description
                'Status'             = $user.Status
                'Disabled'           = $user.Disabled
                'PasswordRequired'   = $user.PasswordRequired
                'PasswordExpires'    = $user.PasswordExpires
                'PasswordChangeable' = $user.PasswordChangeable
                'Lockout'            = $user.Lockout
            }
        }
    }
}

$localUsersDetail | ConvertTo-Json

