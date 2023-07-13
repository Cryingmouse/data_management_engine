param (
    [string]$DirectoryPaths
)

function Get-DirectoryAttributes {
    param (
        [Parameter(Mandatory = $true)]
        [int]$Attributes
    )

    # Attribute bit masks
    $readOnlyMask = 1
    $hiddenMask = 2
    $systemMask = 4
    $directoryMask = 16
    $archiveMask = 32
    $deviceMask = 64
    $normalMask = 128
    $temporaryMask = 256
    $sparseFileMask = 512
    $reparsePointMask = 1024
    $compressedMask = 2048
    $offlineMask = 4096
    $notContentIndexedMask = 8192

    # Parse attributes
    $readOnly = ($Attributes -band $readOnlyMask) -ne 0
    $hidden = ($Attributes -band $hiddenMask) -ne 0
    $system = ($Attributes -band $systemMask) -ne 0
    $directory = ($Attributes -band $directoryMask) -ne 0
    $archive = ($Attributes -band $archiveMask) -ne 0
    $device = ($Attributes -band $deviceMask) -ne 0
    $normal = ($Attributes -band $normalMask) -ne 0
    $temporary = ($Attributes -band $temporaryMask) -ne 0
    $sparseFile = ($Attributes -band $sparseFileMask) -ne 0
    $reparsePoint = ($Attributes -band $reparsePointMask) -ne 0
    $compressed = ($Attributes -band $compressedMask) -ne 0
    $offline = ($Attributes -band $offlineMask) -ne 0
    $notContentIndexed = ($Attributes -band $notContentIndexedMask) -ne 0

    # Return parsed attributes
    return @{
        ReadOnly          = $readOnly
        Hidden            = $hidden
        System            = $system
        Directory         = $directory
        Archive           = $archive
        Device            = $device
        Normal            = $normal
        Temporary         = $temporary
        SparseFile        = $sparseFile
        ReparsePoint      = $reparsePoint
        Compressed        = $compressed
        Offline           = $offline
        NotContentIndexed = $notContentIndexed
    }
}

function Get-DirectoryDetails {
    param (
        [Parameter(Mandatory = $true)]
        [string]$DirectoryPaths
    )

    $directoryPathsArray = $DirectoryPaths -split ';'
    $directoryDetails = @{}

    foreach ($directoryPath in $directoryPathsArray) {
        if (Test-Path -Path $directoryPath -PathType Container) {
            $directory = Get-Item -Path $directoryPath | Select-Object Name, FullName, CreationTime, LastWriteTime, LastAccessTime, Exists, Attributes, Root, Parent
            $attributesDetails = Get-DirectoryAttributes -Attributes $directory.Attributes

            $directoryDetails[$directory.Name] = @{
                "Name"           = $directory.Name
                "FullPath"       = $directory.FullName
                "CreationTime"   = $directory.CreationTime.DateTime
                "LastWriteTime"  = $directory.LastWriteTime.DateTime
                "LastAccessTime" = $directory.LastAccessTime.DateTime
                "Exist"          = $directory.Exists
                "Attributes"     = $attributesDetails
                "ParentFullPath" = $directory.Parent.FullName
            }
        }
    }

    $directoryDetails | ConvertTo-Json
}

Get-DirectoryDetails -DirectoryPaths $DirectoryPaths
