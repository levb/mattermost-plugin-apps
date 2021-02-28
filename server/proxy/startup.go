// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

// <><> const callOnceKey = "CallOnce_key"

// func (adm *Admin) callOnce(f func() error) error {
// 	// Delete previous job
// 	if err := adm.mm.KV.Delete(callOnceKey); err != nil {
// 		return errors.Wrap(err, "can't delete key")
// 	}
// 	// Ensure all instances run this
// 	time.Sleep(10 * time.Second)

// 	adm.mutex.Lock()
// 	defer adm.mutex.Unlock()
// 	value := 0
// 	if err := adm.mm.KV.Get(callOnceKey, &value); err != nil {
// 		return err
// 	}
// 	if value != 0 {
// 		// job is already run by other instance
// 		return nil
// 	}

// 	// job is should be run by this instance
// 	if err := f(); err != nil {
// 		return errors.Wrap(err, "can't run the job")
// 	}
// 	value = 1
// 	ok, err := adm.mm.KV.Set(callOnceKey, value)
// 	if err != nil {
// 		return errors.Wrapf(err, "can't set key %s to %d", callOnceKey, value)
// 	}
// 	if !ok {
// 		return errors.Errorf("can't set key %s to %d", callOnceKey, value)
// 	}
// 	return nil
// }

// func (adm *Admin) expandedCall(sessionToken string, app *apps.App, call *apps.Call, values map[string]string) error {
// 	if call == nil {
// 		return nil
// 	}

// 	if call.Values == nil {
// 		call.Values = map[string]interface{}{}
// 	}
// 	call.Values[apps.PropOAuth2ClientSecret] = app.OAuth2ClientSecret
// 	for k, v := range values {
// 		call.Values[k] = v
// 	}

// 	if call.Expand == nil {
// 		call.Expand = &apps.Expand{}
// 	}
// 	call.Expand.App = apps.ExpandAll
// 	call.Expand.AdminAccessToken = apps.ExpandAll

// 	_, resp := adm.proxy.Call(sessionToken, call)
// 	if resp.Type == apps.CallResponseTypeError {
// 		return errors.Wrapf(resp, "call %s failed", call.URL)
// 	}
// 	return nil
// }
