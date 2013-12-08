import pyparsing

## get_num("1 1/2")
## Parses a string like "1 1/2" into a floating point number like 1.5.
def get_num(string):
	# Parse the number, which may include a fraction.		
	number = pyparsing.Regex(r"\d+(\.\d*)?").setParseAction(lambda t: float(t[0]))
	fraction = number("numerator") + "/" + number("denominator")
	fraction.setParseAction(lambda t: t.numerator / t.denominator)
		
	fractExpr = fraction | number + pyparsing.Optional(fraction)
	fractExpr.setParseAction(lambda t: sum(t))
		
	# Floating representation of the amount.
	try:
		return float(fractExpr.parseString( string ).asList()[0])
	# Exception is thrown if there is no number in the string
	# i.e. ["confectioners' sugar for dusting]
	except pyparsing.ParseException:
		return 0

## get_time("1 day 2 hrs 25 min")
## Parses a string like the example above and converts it into minutes.
def get_time(string):
	tokens = string.split()
	if len(tokens) % 2 != 0:
		raise Exception("Unknown time string format: %s" % string)
	
	units = [unit for i, unit in enumerate(tokens) if (i + 1) % 2 == 0]
	times = [int(time) for i, time in enumerate(tokens) if i % 2 == 0]
	
	minutes = 0
	
	for i, unit in enumerate(units):
		if unit in ('min', 'mins'):
			minutes += times[i]
		elif unit in ('hr', 'hrs'):
			minutes += times[i] * 60
		elif unit in ('day', 'days'):
			minutes += times[i] * 1440
		else:
			raise Exception("Unknown time unit: %s" % unit)
	
	return minutes

def get_servings(string):
	tokens = string.split()
	
	count = 0
	if tokens[1] in ('dozen',):
		count = get_num(tokens[0]) * 12
	else:
		count = get_num(tokens[0])
		
	return count
		
