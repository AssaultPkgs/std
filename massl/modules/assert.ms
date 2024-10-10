
def assertEqual(expected, value)
	if expected == value
		return
	end
	
	throw "assert equal failed"
end

def assertNotEqual(expected, value)
	if expected != value
		return
	end
	
	throw "assert notEqual failed"
end

def assertTrue(value)
	if value
		return
	end
	
	throw "assert true failed"
end

def assertFalse(value)
	if value
		throw "assert false failed"
	end

	return	
end

{
	equal = assertEqual,
	notEqual = assertNotEqual,
	false = assertFalse,
	true = assertTrue
}

