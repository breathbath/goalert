package smoke

import (
	"testing"
	"time"

	"github.com/breathbath/goalert/test/smoke/harness"
)

// TestRotation_Wrap checks that rotations wrap & repeat
func TestRotation_Wrap(t *testing.T) {
	t.Parallel()

	const sql = `
	insert into users (id, name, email)
	values
		({{uuid "uid1"}}, 'bob', 'joe'),
		({{uuid "uid2"}}, 'ben', 'frank');

	insert into user_contact_methods (id, user_id, name, type, value)
	values
		({{uuid "cm1"}}, {{uuid "uid1"}}, 'personal', 'SMS', {{phone "1"}}),
		({{uuid "cm2"}}, {{uuid "uid2"}}, 'personal', 'SMS', {{phone "2"}});

	insert into user_notification_rules (user_id, contact_method_id, delay_minutes)
	values
		({{uuid "uid1"}}, {{uuid "cm1"}}, 0),
		({{uuid "uid2"}}, {{uuid "cm2"}}, 0);

	insert into escalation_policies (id, name)
	values
		({{uuid "eid"}}, 'esc policy');
	insert into escalation_policy_steps (id, escalation_policy_id)
	values
		({{uuid "esid"}}, {{uuid "eid"}});

	insert into schedules (id, name, time_zone)
	values
		({{uuid "sched1"}}, 'default', 'America/Chicago');

	insert into rotations (id, schedule_id, name, type, start_time, shift_length)
	values
		({{uuid "rot1"}}, {{uuid "sched1"}}, 'default rotation', 'hourly', now() - '4 hours'::interval + '10 minute'::interval, 2);

	insert into rotation_participants (rotation_id, user_id, position)
	values
		({{uuid "rot1"}}, {{uuid "uid1"}}, 0),
		({{uuid "rot1"}}, {{uuid "uid2"}}, 1);

	insert into escalation_policy_actions (escalation_policy_step_id, schedule_id)
	values
		({{uuid "esid"}}, {{uuid "sched1"}});

	insert into services (id, escalation_policy_id, name) values
		({{uuid "sid"}}, {{uuid "eid"}}, 'service');

	insert into alerts (service_id, description) values
		({{uuid "sid"}}, 'testing');

	`
	h := harness.NewHarness(t, sql, "ids-to-uuids")
	defer h.Close()

	// with an hourly shift_length of 2, and 2 users; 3hr59min should be user 2, a minute later should be user 1
	sid := h.UUID("sid")
	uid1 := h.UUID("uid1")
	uid2 := h.UUID("uid2")

	h.WaitAndAssertOnCallUsers(sid, uid2)

	h.FastForward(20 * time.Minute)

	h.WaitAndAssertOnCallUsers(sid, uid1)
}
