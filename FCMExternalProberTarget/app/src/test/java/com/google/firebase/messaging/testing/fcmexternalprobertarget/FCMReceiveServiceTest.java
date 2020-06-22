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
import android.text.TextUtils;
import android.util.Log;

import com.google.firebase.messaging.RemoteMessage;

import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.io.IOException;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneId;
import java.util.Collections;
import java.util.Map;
import java.util.Scanner;

import static org.junit.Assert.*;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.fail;
import static org.mockito.Matchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@RunWith (PowerMockRunner.class)
@PrepareForTest({Log.class, RemoteMessage.class})
public class FCMReceiveServiceTest {

    public final String TEST_TOKEN = "TEST_TOKEN";
    public final String SEND_TIME_1 = "0200";
    public final String SEND_TIME_2 = "0300";
    public Clock testClock;

    @Rule
    public TemporaryFolder testFolder = new TemporaryFolder();

    @Mock
    Context mockContext = mock(Context.class);

    @Before
    public void initTests() {
        testClock = Clock.fixed(Instant.EPOCH, ZoneId.of("UTC"));
        PowerMockito.mockStatic(Log.class);
        PowerMockito.mockStatic(TextUtils.class);
    }

    @Test
    public void onNewTokenTest_expected() {
        Clock testClock = Clock.fixed(Instant.EPOCH, ZoneId.of("UTC"));
        FCMReceiveService service = new FCMReceiveService(mockContext, true, testClock);
        try {
            File validDirectory = testFolder.newFolder();
            when(mockContext.getExternalFilesDir(anyString())).thenReturn(validDirectory);

            service.onNewToken(TEST_TOKEN);

            Scanner scanner = new Scanner(new File(validDirectory, "token.txt"));
            assertEquals(TEST_TOKEN, scanner.nextLine());
            assertFalse(scanner.hasNext());
        } catch (IOException exception){
            fail(exception.toString());
        }
    }

    @Test
    public void onMessageReceivedTest_expected() {
        FCMReceiveService service = new FCMReceiveService(mockContext, true, testClock);
        try {
            File validDirectory = testFolder.newFolder();
            RemoteMessage testMessage = PowerMockito.mock(RemoteMessage.class);
            Map<String,String> testData = Collections.singletonMap("sendTime", SEND_TIME_1);

            PowerMockito.when(testMessage.getData()).thenReturn(testData);
            when(mockContext.getExternalFilesDir(anyString())).thenReturn(validDirectory);

            service.onMessageReceived(testMessage);

            Scanner scanner = new Scanner(new File(validDirectory, "logs/" + SEND_TIME_1 + ".txt"));
            assertEquals(testClock.instant().getEpochSecond(), scanner.nextLong());
            assertFalse(scanner.hasNext());
        } catch (IOException exception) {
            fail(exception.toString());
        }
    }
    @Test
    public void onMessageReceivedTest_twoMessages() {
        FCMReceiveService service = new FCMReceiveService(mockContext, true, testClock);
        try {
            File validDirectory = testFolder.newFolder();
            RemoteMessage testMessage = PowerMockito.mock(RemoteMessage.class);
            Map<String,String> testData = Collections.singletonMap("sendTime", SEND_TIME_1);
            RemoteMessage testMessage2 = PowerMockito.mock(RemoteMessage.class);
            Map<String,String> testData2 = Collections.singletonMap("sendTime", SEND_TIME_2);

            PowerMockito.when(testMessage.getData()).thenReturn(testData);
            PowerMockito.when(testMessage2.getData()).thenReturn(testData2);
            when(mockContext.getExternalFilesDir(anyString())).thenReturn(validDirectory);

            service.onMessageReceived(testMessage);
            service.onMessageReceived(testMessage2);

            Scanner scanner = new Scanner(new File(validDirectory, "logs/" + SEND_TIME_1 + ".txt"));
            assertEquals(testClock.instant().getEpochSecond(), scanner.nextLong());
            assertFalse(scanner.hasNext());

            scanner = new Scanner(new File(validDirectory, "logs/" + SEND_TIME_2 + ".txt"));
            assertEquals(testClock.instant().getEpochSecond(), scanner.nextLong());
            assertFalse(scanner.hasNext());
        } catch (IOException exception) {
            fail(exception.toString());
        }
    }
}
