/*
    Cinnamon UCI chess engine
    Copyright (C) Giuseppe Cannella

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#pragma once

#include <thread>
#include <mutex>
#include "ObserverThread.h"
#include "../namespaces.h"
#include "CoutSync.h"
#include "ConditionVariable.h"

using namespace std;

class Runnable {
public:
    virtual void run() = 0;
};

class Thread : virtual public Runnable {

private:
    bool running = true;
    int threadID = -1;
    ObserverThread *observer = nullptr;
    ConditionVariable cv;
    thread theThread;

    Runnable *execRunnable;

    static void *__run(void *cthis) {
        static_cast<Runnable *>(cthis)->run();
        static_cast<Thread *>(cthis)->notifyEndThread((static_cast<Thread *>(cthis))->getId());

        return nullptr;
    }

public:

    Thread() {
        execRunnable = this;
    }

    void registerObserverThread(ObserverThread *obs) {
        observer = obs;
    }

    void notifyEndThread(int threadID) {
        if (observer != nullptr) {
            observer->observerEndThread(threadID);
        }
    }

    virtual ~Thread() {

    }

    void checkWait() {
        while (!running) {
            cv.wait();
        }
    }

    void notify() {
        cv.notify_all();
    }

    void start() {
        join();
        theThread = thread(__run, execRunnable);
    }

    void join() {
        if (theThread.joinable()) {
#ifdef DEBUG_MODE
            CoutSync() << "join: " << threadID;
#endif
            theThread.join();
        }
    }

//    bool isJoinable() {
//        return  theThread.joinable();
//    }

    void detach() {
        theThread.detach();
    }

    int getId() const {
        return threadID;
    }

    void setId(int threadID) {
        Thread::threadID = threadID;
    }

    void sleep(bool b) {
        running = !b;
    }

//    void stop() {
//        if (theThread) {
//             theThread.detach();
//            delete theThread;
//            theThread = nullptr;
//        }
//    }
};
