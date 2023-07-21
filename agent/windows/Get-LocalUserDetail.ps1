param (
    [String] $UserNames
)

$localUsersDetail = @()

$UserNamesArray = $UserNames -split ','

$allUsers = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True'"

if (-not $UserNames) {
    foreach ($user in $allUsers) {
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
else {
    $users = $allUsers | Where-Object { $_.Name -in $UserNamesArray }
    foreach ( $user in $users ) {
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

$localUsersDetail | ConvertTo-Json

