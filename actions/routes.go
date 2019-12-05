package actions

var Routes []route = []route{
	{
		path:   "login",
		action: Login,
	},
	{
		path:   "invite",
		action: Invite,
	},
	{
		path:   "add",
		action: Add,
	},
	{
		path:   "run",
		action: Run,
	},
	{
		path:   "test",
		action: Test,
	},
	{
		path:   "accept",
		action: Accept,
	},
	{
		path:   "composer",
		action: Composer,
	},
	{
		path:   "allusers",
		action: Allusers,
	},
	{
		path:   "alllanguages",
		action: Alllanguages,
	},
}
