/*
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.google.firebase.messaging.testing.fcmexternalprobertarget;

import android.content.Context;
import android.content.Intent;
import android.widget.TextView;

import androidx.test.platform.app.InstrumentationRegistry;

import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class LogReceiverInstrumentedTest {

    public final String TEST_STRING = "TEST_STRING";
    public Context testContext;
    public TextView testView;
    public Intent testIntent;

    @Before
    public void init() {
        testContext = InstrumentationRegistry.getInstrumentation().getTargetContext();
        testView = new TextView(testContext);
        testIntent = new Intent();
    }

    @Test
    public void onReceiveTest_expected() {
        testIntent.putExtra("logText", TEST_STRING);
        LogReceiver testReceiver = new LogReceiver(testView);

        testReceiver.onReceive(testContext, testIntent);

        assertEquals(TEST_STRING + "\n", testView.getText().toString());
    }

    @Test
    public void onReceiveTest_null() {
        testIntent.putExtra("logText", (String) null);
        LogReceiver testReceiver = new LogReceiver(testView);

        testReceiver.onReceive(null, testIntent);

        assertEquals("Log Error: No log text supplied\n", testView.getText().toString());
    }

    @Test
    public void onReceiveTest_twice() {
        testIntent.putExtra("logText", TEST_STRING);
        LogReceiver testReceiver = new LogReceiver(testView);

        testReceiver.onReceive(testContext, testIntent);
        testReceiver.onReceive(testContext, testIntent);

        assertEquals(TEST_STRING + "\n" + TEST_STRING + "\n", testView.getText().toString());
    }
}
