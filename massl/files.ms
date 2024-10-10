
dirinfo = new System.IO.DirectoryInfo(".")

for fileinfo in dirinfo.GetFiles()
	println(fileinfo.Name)
end
