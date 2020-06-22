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

import android.content.Intent;
import android.widget.TextView;

import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.runners.MockitoJUnitRunner;

import java.util.Collections;
import java.util.Map;

import static org.junit.Assert.*;

@RunWith(MockitoJUnitRunner.class)
public class LogReceiverTest {

    public static class MockTextView extends TextView {
        private StringBuilder textBuffer;

        public MockTextView(StringBuilder textBuffer) {
            super(null);
            this.textBuffer = textBuffer;
        }

        @Override
        public void append(CharSequence toAdd, int start, int end) {
            textBuffer.append(toAdd);
        }

        @Override
        public int length() {
            return textBuffer.length();
        }

        public String getContents() {
            return textBuffer.toString();
        }
    }

    public static class MockIntent extends Intent {
        private Map<String, String> extraData;

        public MockIntent(String extraDataKey, String extraDataValue) {
            extraData = Collections.singletonMap(extraDataKey, extraDataValue);
        }

        @Override
        public String getStringExtra(String extraDataKey) {
            return extraData.get(extraDataKey);
        }
    }

    public final String TEST_STRING = "TEST_STRING";

    @Rule
    MockTextView mockedTextView = new MockTextView(new StringBuilder());

    @Rule
    LogReceiver testReceiver = new LogReceiver(mockedTextView);

    @Test
    public void onReceiveTest_expected() {
        MockIntent testIntent = new MockIntent("logText", TEST_STRING);

        testReceiver.onReceive(null, testIntent);

        assertEquals(TEST_STRING + "\n", mockedTextView.getContents());
    }

    @Test
    public void onReceiveTest_null() {
        MockIntent testIntent = new MockIntent("logText", null);

        testReceiver.onReceive(null, testIntent);

        assertEquals("Log Error: No log text supplied\n", mockedTextView.getContents());
    }
}


