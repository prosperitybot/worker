package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// DiscordEpoch is the constant in time.Duration (nanoseconds)
// since Unix epoch.
const DiscordEpoch = 1420070400000 * time.Millisecond

// DurationSinceEpoch returns the duration from the Discord epoch to current.
func DurationSinceEpoch(t time.Time) time.Duration {
	return time.Duration(t.UnixNano()) - DiscordEpoch
}

type Snowflake int64

// NullSnowflake gets encoded into a null. This is used for
// optional and nullable snowflake fields.
const NullSnowflake = ^Snowflake(0)

func NewSnowflake(t time.Time) Snowflake {
	return Snowflake((DurationSinceEpoch(t) / time.Millisecond) << 22)
}

func (sf Snowflake) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf(`"%d"`, sf)
	return []byte(jsonValue), nil
}

func (sf *Snowflake) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return errors.New("error unmarshalling JSON")
	}

	sfNum, err := strconv.ParseUint(unquotedJSONValue, 10, 64)
	if err != nil {
		return errors.New("unsupported JSON value")
	}

	*sf = Snowflake(sfNum)

	return nil
}
