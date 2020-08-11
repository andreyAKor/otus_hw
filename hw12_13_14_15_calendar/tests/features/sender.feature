# file: features/sender.feature

Feature: Sender service in work
	Scenario: Event is received
		When I send "POST" request to "http://calendar:5080/create" with "application/x-www-form-urlencoded" create new event, with data:
		"""
		title=title-456
		date=@timeNow
		duration=1h
		descr=descr-789
		user_id=1
		duration_start=0
		"""
		Then I receive event with data:
		"""
		OK
		"""
