/*
    https://github.com/gekomad/BlockingThreadPool
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
#include "../namespaces/def.h"
#include <condition_variable>
#include "ThreadBase.h"

using namespace std;


class ThreadReuse : public ThreadBase {

private:
    bool running = true;
    int threadID;
    ObserverThread *observer = nullptr;
    condition_variable cv;
    thread theThread;
    bool started = false;

    Runnable *execRunnable;
protected:
    static void *__run(void *cthis) {
        while (true) {
            static_cast<Runnable *>(cthis)->run();

            static_cast<Runnable *>(cthis)->endRun();

            static_cast<ThreadReuse *>(cthis)->notifyEndThread((static_cast<ThreadReuse *>(cthis))->getId());

        }
    }

public:
    bool reuse;

    ThreadReuse() {
        execRunnable = this;
    }

    void registerObserverThread(ObserverThread *obs) {
        observer = obs;
    }

    void notifyEndThread(int i) {
        ASSERT(observer);
        if (observer != nullptr) {
            observer->observerEndThread(i);
            mutex mtx;

            unique_lock<mutex> lck(mtx);
            cv.wait(lck);

        }
    }

    virtual ~ThreadReuse() {
        join();
    }

    void checkWait() {
        while (!running) {
            mutex mtx;
            unique_lock<mutex> lck(mtx);
            cv.wait(lck);
        }
    }

    void notify() {
        cv.notify_all();
    }

    void start() {
        if (!started) {
            started = true;
            ASSERT(!isJoinable());
            theThread = thread(__run, execRunnable);
        } else {
            cv.notify_all();
        }
    }

    void join() {

    }

    void detach() {
        theThread.detach();
    }

    int getId() const {
        return threadID;
    }

    void setId(int id) {
        threadID = id;
    }

    void threadSleep(bool b) {
        running = !b;
    }

    bool isJoinable() {
        return theThread.joinable();
    }

    void setSleep(bool b) {
        running = !b;
    }

};
