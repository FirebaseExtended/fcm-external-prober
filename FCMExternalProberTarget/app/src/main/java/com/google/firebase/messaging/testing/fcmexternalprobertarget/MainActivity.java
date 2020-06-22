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

import androidx.appcompat.app.AppCompatActivity;
import androidx.localbroadcastmanager.content.LocalBroadcastManager;

import android.content.IntentFilter;
import android.os.Bundle;
import android.widget.TextView;

/**
 * App that receives messages from Firebase Cloud Messaging and logs reception
 */
public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        createBroadcastReceiver();
    }

    private void createBroadcastReceiver() {
       TextView logText = findViewById(R.id.logTextView);
        LogReceiver logReceiver = new LogReceiver(logText);
        IntentFilter logFilter = new IntentFilter("updateUI");

        LocalBroadcastManager.getInstance(this).registerReceiver(logReceiver, logFilter);
    }
}