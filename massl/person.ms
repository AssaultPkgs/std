
class Person
	define initialize(firstname, lastname)
		self.firstname = firstname
		self.lastname = lastname
	end
	
	define getName()
		return self.lastname + ", " + self.firstname
	end
end

adam = new Person("Adam", "TheFirst")

println(adam.getName())
