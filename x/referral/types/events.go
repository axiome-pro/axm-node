package types

func (EventStatusUpdated) XXX_MessageName() string { return "status_updated" }

func (EventStatusWillBeDowngraded) XXX_MessageName() string { return "status_will_be_downgraded" }

func (EventStatusDowngradeCanceled) XXX_MessageName() string { return "status_downgrade_canceled" }
