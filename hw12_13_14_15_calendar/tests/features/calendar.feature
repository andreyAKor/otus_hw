# file: features/calendar.feature

Feature: Calendar service in work
	Scenario: Succesful response for list by date
		When I send "GET" request to "http://calendar:5080/list/date?date=2018-09-22"
		Then The response code should be 200

	Scenario: Bad request for list by date
		When I send "GET" request to "http://calendar:5080/list/date?date=***-++-22"
		Then The response code should be 400

	Scenario: Succesful response for list by week
		When I send "GET" request to "http://calendar:5080/list/week?start=2018-09-22"
		Then The response code should be 200

	Scenario: Bad request for list by week
		When I send "GET" request to "http://calendar:5080/list/week?start=***-++-22"
		Then The response code should be 400

	Scenario: Succesful response for list by month
		When I send "GET" request to "http://calendar:5080/list/month?start=2018-09-22"
		Then The response code should be 200

	Scenario: Bad request for list by month
		When I send "GET" request to "http://calendar:5080/list/month?start=***-++-22"
		Then The response code should be 400

	Scenario: Succesful create new event
		When I send "POST" request to "http://calendar:5080/create" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2020-08-02T16:02:31+03:00
		duration=1h
		descr=descr-456
		user_id=32
		duration_start=2h
		"""
		Then The response code should be 201
		And I receive response with data:
		"""
		{
			"data": {
				"id": 1
			}
		}
		"""

	Scenario: Succesful create new additional event
		When I send "POST" request to "http://calendar:5080/create" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2020-08-03T16:02:31+03:00
		duration=1h
		descr=descr-456
		user_id=32
		duration_start=2h
		"""
		Then The response code should be 201
		And I receive response with data:
		"""
		{
			"data": {
				"id": 2
			}
		}
		"""

	Scenario: Create dublicate event
		When I send "POST" request to "http://calendar:5080/create" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2020-08-02T16:02:31+03:00
		duration=1h
		descr=descr-456
		user_id=32
		duration_start=2h
		"""
		Then The response code should be 409
		And I receive response with data:
		"""
		{
			"error": "can't create new event: event at this time is busy"
		}
		"""

	Scenario: Bad request for creating new event
		When I send "POST" request to "http://calendar:5080/create" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=****-++-02T16:02:31+03:00
		duration=*h
		descr=descr-456
		user_id=32
		duration_start=*h
		"""
		Then The response code should be 400

	Scenario: Succesful update current event
		When I send "PUT" request to "http://calendar:5080/update?id=1" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2018-09-22T12:42:31+03:00
		duration=0h
		descr=descr-456
		user_id=32
		duration_start=0
		"""
		Then The response code should be 200
		And I receive response with data:
		"""
		{}
		"""

	Scenario: Bad request for updating event in case 1
		When I send "PUT" request to "http://calendar:5080/update?id=1" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=****-++-22T12:42:31+03:00
		duration=0h
		descr=descr-456
		user_id=32
		duration_start=0
		"""
		Then The response code should be 400

	Scenario: Bad request for updating event in case 2
		When I send "PUT" request to "http://calendar:5080/update?id=**" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2018-09-22T12:42:31+03:00
		duration=0h
		descr=descr-456
		user_id=32
		duration_start=0
		"""
		Then The response code should be 400

	Scenario: Not found event for updating
		When I send "PUT" request to "http://calendar:5080/update?id=999999999" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2018-09-22T12:42:31+03:00
		duration=0h
		descr=descr-456
		user_id=32
		duration_start=0
		"""
		Then The response code should be 404

	Scenario: Dublicate event data for updating
		When I send "PUT" request to "http://calendar:5080/update?id=1" with "application/x-www-form-urlencoded" data:
		"""
		title=title-123
		date=2020-08-03T16:02:31+03:00
		duration=1h
		descr=descr-456
		user_id=32
		duration_start=2h
		"""
		Then The response code should be 409

	Scenario: Succesful deleting current event
		When I send "DELETE" request to "http://calendar:5080/delete?id=1"
		Then The response code should be 200
		And I receive response with data:
		"""
		{}
		"""

	Scenario: Bad request for deleting event
		When I send "DELETE" request to "http://calendar:5080/delete?id=*"
		Then The response code should be 400

	Scenario: Not found event for deleting
		When I send "DELETE" request to "http://calendar:5080/delete?id=999999999"
		Then The response code should be 404
