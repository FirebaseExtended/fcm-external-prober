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

import androidx.test.platform.app.InstrumentationRegistry;

import com.google.firebase.messaging.RemoteMessage;

import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

import java.io.File;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneId;
import java.util.Scanner;

public class FCMReceiveServiceInstrumentedTest {

    public final String TEST_TOKEN = "TEST_TOKEN";
    public final String SEND_TIME = "0200";
    public final String TYPE = "0";
    public Context testContext;
    public Clock testClock;
    public FCMReceiveService service;

    @Before
    public void init() {
        testContext = InstrumentationRegistry.getInstrumentation().getTargetContext();
        testClock = Clock.fixed(Instant.EPOCH, ZoneId.of("UTC"));
        service = new FCMReceiveService(testContext, true, testClock);
    }

    @Test
    public void onNewTokenTest_expected() throws Exception {
        service.onNewToken(TEST_TOKEN);

        File testFile = new File(testContext.getExternalFilesDir(null), "token.txt");
        assertTrue(testFile.exists());
        Scanner scanner = new Scanner(testFile);
        assertEquals(TEST_TOKEN, scanner.nextLine());
        assertFalse(scanner.hasNext());
    }

    @Test
    public void onMessageReceivedTest_expected() throws Exception {
        RemoteMessage testMessage = new RemoteMessage.Builder("TEST")
                .addData("sendTime", SEND_TIME).addData("type", TYPE).build();

        service.onMessageReceived(testMessage);

        File testFile = new File(testContext.getExternalFilesDir(null),
                "logs/" + TYPE + SEND_TIME + ".txt");
        assertTrue(testFile.exists());
        Scanner scanner = new Scanner(testFile);
        assertEquals(testClock.instant().getEpochSecond(), scanner.nextLong());
        assertFalse(scanner.hasNext());
    }
}
