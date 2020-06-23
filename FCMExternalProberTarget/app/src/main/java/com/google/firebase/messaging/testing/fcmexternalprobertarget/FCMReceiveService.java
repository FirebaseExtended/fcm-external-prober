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
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.annotation.VisibleForTesting;
import androidx.localbroadcastmanager.content.LocalBroadcastManager;

import com.google.firebase.messaging.FirebaseMessagingService;
import com.google.firebase.messaging.RemoteMessage;

import java.time.Clock;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;

/**
 * Service that handles interaction with Firebase Cloud Messaging, and logs the results of the
 * messages both to the UI and files.
 */
public class FCMReceiveService extends FirebaseMessagingService {

    private Context context;
    private boolean viewLogging;
    private Clock logTimer;

    public FCMReceiveService () {
        this.context = this;
        viewLogging = false;
        logTimer = Clock.systemUTC();
    }

    /**
     * Create an instance for testing
     * @param context Mocked Context object
     * @param viewLogging Whether logs should be written to UI
     */
    @VisibleForTesting
    public FCMReceiveService (Context context, boolean viewLogging, Clock logTimer) {
        this.context = context;
        this.viewLogging = viewLogging;
        this.logTimer = logTimer;
    }

    @Override
    public void onNewToken(@NonNull String token) {
        // Store token in a file on external storage so it can be accessed by the probe
        File tokenFile = makeExternalFile("","token.txt");
        writeToFile(tokenFile, token);
        logToUI("Info", "Token generated");
    }

    @Override
    public void onMessageReceived(@NonNull RemoteMessage remoteMessage) {
        long receivedTime = logTimer.instant().getEpochSecond();
        String sendTime = remoteMessage.getData().get("sendTime");
        File logFile = makeExternalFile("logs", sendTime + ".txt");
        writeToFile(logFile, Long.toString(receivedTime));
        logToUI("Info","Message Received: " + sendTime + ".txt");
    }


    private void logToUI(String tag, String logText) {
        Log.d(tag, logText);
        if (viewLogging) {
            return;
        }
        Intent intent = new Intent("updateUI");
        intent.putExtra("logText", tag + ": " + logText);
        LocalBroadcastManager.getInstance(this).sendBroadcast(intent);
    }

    private File makeExternalFile(String path, String fileName) {
        File newPath = new File(context.getExternalFilesDir(null), path);
        if (!newPath.exists()) {
            if (!newPath.mkdirs()) {
                logToUI("Error", "Unable to create directory " + path);
            }
            else {
                logToUI("Info", "Directory " + path + " created");
            }
        }
        return new File(newPath, fileName);
    }

    private void writeToFile(File writeFile, String writeText) {
        try {
            FileWriter outputWriter = new FileWriter(writeFile);
            outputWriter.write(writeText, 0, writeText.length());
            outputWriter.close();
        } catch (IOException exception) {
            logToUI("Error", exception.toString());
        }
    }
}
