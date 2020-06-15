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

package com.google.fcmexternalprobertarget;

import android.content.Intent;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.localbroadcastmanager.content.LocalBroadcastManager;

import com.google.firebase.messaging.FirebaseMessagingService;
import com.google.firebase.messaging.RemoteMessage;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;

/**
 * Service that handles interaction with Firebase Cloud Messaging, and logs the results of the
 * messages both to the UI and files.
 */
public class FCMReceiveService extends FirebaseMessagingService {

    @Override
    public void onNewToken(@NonNull String token) {
        // Store token in a file on external storage so it can be accessed by the probe
        File tokenFile = makeExternalFile("token.txt");
        writeToExternalFile(tokenFile, token);
        logToUI("Info", "Token generated");
    }

    @Override
    public void onMessageReceived(@NonNull RemoteMessage remoteMessage) {
        long receivedTime = System.currentTimeMillis();
        String sentTime = remoteMessage.getData().get("sendTime");
        File logFile = makeExternalFile(sentTime + ".txt");
        writeToExternalFile(logFile, Long.toString(receivedTime));
        logToUI("Info" ,"Message Received");
    }

    private void logToUI(String tag, String logText) {
        Log.d(tag, logText);
        Intent intent = new Intent("updateUI");
        intent.putExtra("logText", tag + ": " + logText);
        LocalBroadcastManager.getInstance(this).sendBroadcast(intent);
    }

    private File makeExternalFile(String path) {
        File newFile = new File(getExternalFilesDir(null), path);
        if (!newFile.mkdirs()) {
            logToUI("Error", "Unable to create specified file directory: " + path);
        }
        return newFile;
    }

    private void writeToExternalFile(File writeFile, String writeText) {
        try {
            FileWriter outputWriter = new FileWriter(writeFile);
            outputWriter.write(writeText, 0, writeText.length());
            outputWriter.close();
        }
        catch (IOException exception) {
            logToUI("Error:", "Unable to write token to file");
        }
    }
}
