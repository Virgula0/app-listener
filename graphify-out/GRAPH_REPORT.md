# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 562 files · ~1,789,089 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7310 nodes · 33404 edges · 253 communities (211 shown, 42 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2995 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8531697e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- Guard
- mkdirs
- EventType
- uninstallAll
- editView
- fscrypt.go
- Network
- Candidate
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- FileTarget
- guard.bpf.c
- LibraryClosure
- RunCmd
- updatePage
- go_pkg_strings
- classifyGroup
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- GuardInodeKey
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- os.FileMode
- RefreshSection
- IntegrationSuite
- absPath
- monitorUnitTest
- cursor.go
- Config
- .resKey
- stdio
- update_test.go
- fadeScroll
- New
- .Call
- Client
- daemonconfig.go
- session
- .apply
- githubRelease
- mutedText
- engine
- Monitor
- Prereq
- .FilterExistingWhitelist
- NetworkMonitor
- .filterEvent
- binaryList
- unix_connect.c
- wm_test.go
- IntegrationSuite
- linkView
- FindDaemonProcesses
- Vault
- daemon_multicall_test.go
- journal
- fcntl
- guardUnitTest
- LaunchRule
- codecheck/main.go
- NetGuardEvent
- User
- process_vm_readv.c
- os.File
- recordFlow
- os.FileInfo
- rcUser
- golang.org/x/sys/unix.Stat_t
- github.com/cilium/ebpf.CollectionSpec
- exe_supersede.h
- loginServer
- Vault
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- removeKeyAndEmptyDir
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- BinaryEntry
- runDaemon
- Well-known bypass classes
- go_pkg_os
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- newTestVetter
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- backend.go
- ConfinedEntry
- installPage
- .judgeProcess
- IntegrationSuite
- maintainerKeeper
- downloadAndVerify
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- UI
- go_pkg_path_filepath
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- classify
- make test-integration (rootful Docker suite)
- newTestTarget
- catalogRefresher
- planUpdaters
- trace-app-libs.sh
- IntegrationSuite
- Finder
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- softenAutomatedRefreshErr
- serve.go
- liveSession
- github.com/charmbracelet/bubbles/textarea.Model
- Ledger
- daemon: permanent protection with encryption at rest
- GuardEvent
- buildTrustedSet
- RunWith
- fakeStore
- FileEvent
- IntegrationSuite
- buildOneGuard
- watchSet
- CopyTreeWithProgress
- sources
- Installation
- backend_test.go
- graphify-refresh.sh
- bufio.Reader
- NetEventType
- runOfflineEdit
- Active
- feed.go
- shell
- sources_test.go
- procstats.go
- Run
- SelectOrphans
- Limitations
- dialLiveSession
- IntegrationSuite
- TrustGuard
- net_tester/main.go
- selection_test.go
- deploy
- IntegrationSuite
- time.Time
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- IntegrationSuite
- testLimiter
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- recoverFileInPlace
- pathCache
- pinstate_test.go
- DaemonEvent
- WithBottomBar
- Recorder
- appTheme
- monitor
- IntegrationSuite
- mustSubkey
- bpfstats/main.go
- .eventTypeAllowed
- KernelDev
- migrate_test.go
- .exeKeys
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- SecretStoreNotes
- parseDaemonEvents
- .run
- WithSuspectsFrom
- Gate
- ValidateTargets
- configedit_test.go
- reloadOnce
- supersede_probe.c
- github.com/charmbracelet/bubbletea.Cmd
- TestCleanRel
- mountTable
- eventUnitTest
- github.com/Virgula0/app-listener/tools/codecheck
- fileedit_symlink_test.go
- procRec
- Question
- bleedLayout
- inode_dev.h
- Troubleshooting
- ReadInstallCert
- GenerateConf
- ResolvePinBase
- shQuote
- taint_victim.c
- parsePasswd
- checkRunningBinaryMatchesInstalled
- .showSized
- rgb
- monitor: see what touches a path
- fileID
- plan_test.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 146 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (253 total, 42 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.18
Nodes (31): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (169): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+161 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): Hash(), configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.16
Nodes (14): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+6 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (86): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+78 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "Guard"
Cohesion: 0.06
Nodes (20): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), binaryVerifyState, ReplacementCheck, TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), Guard (+12 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "EventType"
Cohesion: 0.10
Nodes (19): selfEventsFor(), BusClient, EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), AttestBusClient(), busClientKey() (+11 more)

### Community 11 - "uninstallAll"
Cohesion: 0.06
Nodes (46): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+38 more)

### Community 12 - "editView"
Cohesion: 0.11
Nodes (8): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), newEditView(), relOf()

### Community 13 - "fscrypt.go"
Cohesion: 0.12
Nodes (17): ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey(), isNotEncryptedErr(), MasterKeyExists() (+9 more)

### Community 14 - "Network"
Cohesion: 0.13
Nodes (7): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes()

### Community 15 - "Candidate"
Cohesion: 0.14
Nodes (21): TestUncoveredCandidates(), uncoveredCandidates(), groupCandidates(), pickDirectories(), pickFromCandidates(), TestGroupCandidates(), Candidate, catalogGroup (+13 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+56 more)

### Community 18 - "navItem"
Cohesion: 0.04
Nodes (19): clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+11 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (19): ownerPaths(), backingDeviceUnion(), buildGuards(), systemPatterns(), fakeDaemon, Resource, BackingDevices(), GuardRepository (+11 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (38): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+30 more)

### Community 24 - "FileTarget"
Cohesion: 0.11
Nodes (13): guardPill(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget(), NormalizeSpec() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (93): inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+85 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.07
Nodes (31): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+23 more)

### Community 27 - "RunCmd"
Cohesion: 0.05
Nodes (61): offerDaemonStart(), installFacultatives(), installSSHAgentFacultative(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), installGUI() (+53 more)

### Community 28 - "updatePage"
Cohesion: 0.15
Nodes (13): installedBinary, latestRelease, stagedBinary, updatePage, hashFile(), hashReader(), readInstalled(), stageRelease() (+5 more)

### Community 29 - "go_pkg_strings"
Cohesion: 0.07
Nodes (9): needsBinary(), TestNeedsBinary(), Execute(), initLogger(), Mode, Request, main(), i3Result (+1 more)

### Community 30 - "classifyGroup"
Cohesion: 0.12
Nodes (23): loaderName(), LoaderRefusal(), soNames(), appletLine(), beyondEvents(), classifyGroup(), classifyLibs(), finishBinary() (+15 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "GuardInodeKey"
Cohesion: 0.06
Nodes (29): inodeBatchFn, refusedReplacements, rootHandle, openRoot(), RegisterBusDir(), appendInodeRow(), batchInodeRows(), inodeRowsOf() (+21 more)

### Community 34 - "picker"
Cohesion: 0.06
Nodes (37): adaptiveGrid, gapVBox, keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch (+29 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.26
Nodes (23): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), get() (+15 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.16
Nodes (10): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 39 - "os.FileMode"
Cohesion: 0.11
Nodes (22): entryOf(), readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned() (+14 more)

### Community 40 - "RefreshSection"
Cohesion: 0.09
Nodes (28): logRefreshChange(), TestLogRefreshChange(), refreshCatalogConfig(), TestPatchCatalogSectionGroupedConfig(), vaultScan(), RefreshOptions, SectionScan, ConfirmOverwrite() (+20 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (9): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath(), TestMain(), helperChild() (+1 more)

### Community 43 - "monitorUnitTest"
Cohesion: 0.15
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 45 - "Config"
Cohesion: 0.05
Nodes (67): newBinaryVetter(), openBinaryVetter(), newCatalogRefresher(), sameConfig(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors() (+59 more)

### Community 46 - ".resKey"
Cohesion: 0.09
Nodes (9): exeRow, supersededBinary, tempGrant, tempMaskOp, tempRow, liveTrust(), planTempBlock(), stripPinnedAllow() (+1 more)

### Community 48 - "update_test.go"
Cohesion: 0.10
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanInstalledAcrossTagFormats() (+9 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (14): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide(), TestFadeBarDragMovesContent() (+6 more)

### Community 50 - "New"
Cohesion: 0.14
Nodes (22): encryptDirectories(), harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode() (+14 more)

### Community 51 - ".Call"
Cohesion: 0.08
Nodes (13): backendLink, forwardView, recordReview, trustPage, Event, Hello, PwState, TrustRow (+5 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (10): Client, launch(), TestHelperBackend(), TestReaderCapsLines(), Reader, Writer, NewReader(), NewWriter() (+2 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.07
Nodes (67): SecretScope, watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective() (+59 more)

### Community 54 - "session"
Cohesion: 0.13
Nodes (8): TestReadFrame(), decodeFrame(), encodeWithSerial(), header(), readFrame(), queued, senderQueue, session

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (29): checkGrid, filesPage, forwardForm, homePage, multiAnswer, protectedPage, recordForm, recordView (+21 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (30): engine, portOwner, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), nearestRootClaims() (+22 more)

### Community 59 - "Monitor"
Cohesion: 0.15
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 60 - "Prereq"
Cohesion: 0.29
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 61 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 62 - "NetworkMonitor"
Cohesion: 0.13
Nodes (8): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 64 - "binaryList"
Cohesion: 0.09
Nodes (15): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+7 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm_test.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "linkView"
Cohesion: 0.31
Nodes (5): linkTarget(), newLinkView(), TestCheckLinksJudgesEachTarget(), entryKind, linkView

### Community 69 - "FindDaemonProcesses"
Cohesion: 0.16
Nodes (13): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), anyDaemonRunning(), DaemonRunning() (+5 more)

### Community 70 - "Vault"
Cohesion: 0.12
Nodes (19): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+11 more)

### Community 71 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 72 - "journal"
Cohesion: 0.14
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "guardUnitTest"
Cohesion: 0.07
Nodes (16): addErrKind, bpfGuardEvent, classifyAddErr(), decodeGuardEvent(), eventMask(), fsGateLabel(), parseGuardEvent(), processGateLabel() (+8 more)

### Community 75 - "LaunchRule"
Cohesion: 0.16
Nodes (14): allowedEnv(), binaryLaunchRules(), TestAllowedEnvKeepsOnlyRootPlaced(), argvEntry, LaunchRule, BuiltinAllowedEnv(), BuiltinLaunchRules(), canonLaunchRules() (+6 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetGuardEvent"
Cohesion: 0.05
Nodes (48): DrainUntilSignal(), BuildEBPFTargets(), CheckBPFLSM(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag() (+40 more)

### Community 78 - "User"
Cohesion: 0.06
Nodes (42): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemRule(), entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry() (+34 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.04
Nodes (72): multicallAdmissible(), scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun() (+64 more)

### Community 81 - "recordFlow"
Cohesion: 0.12
Nodes (22): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+14 more)

### Community 82 - "os.FileInfo"
Cohesion: 0.17
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 83 - "rcUser"
Cohesion: 0.13
Nodes (40): TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv() (+32 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.11
Nodes (23): placedWalk, treeJudge, confineBelow(), inUserHome(), PathUnder(), vouchedFd(), GlobSystemPlaced(), placedDir() (+15 more)

### Community 85 - "github.com/cilium/ebpf.CollectionSpec"
Cohesion: 0.23
Nodes (12): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+4 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 87 - "loginServer"
Cohesion: 0.08
Nodes (22): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+14 more)

### Community 88 - "Vault"
Cohesion: 0.11
Nodes (19): Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode(), inodeOf(), mustRead(), mustWrite() (+11 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.06
Nodes (13): daemonPage, logGrid, logHook, logLine, runHost, runView, journalLine(), TestJournalLine() (+5 more)

### Community 92 - "removeKeyAndEmptyDir"
Cohesion: 0.67
Nodes (3): removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir()

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.17
Nodes (6): eventsetSummary(), eventTypeKey(), modeLabel(), NewNetGuard(), hook, NetGuard

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "BinaryEntry"
Cohesion: 0.07
Nodes (28): deferredBinary, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), canonicalPaths() (+20 more)

### Community 101 - "runDaemon"
Cohesion: 0.14
Nodes (13): ServeConfig, newControlManager(), startControlManager(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI(), runServedTUI() (+5 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "go_pkg_os"
Cohesion: 0.07
Nodes (20): releaseBuild(), Run(), TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), asRoot(), runInstall(), closeDumpFile() (+12 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "newTestVetter"
Cohesion: 0.13
Nodes (23): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+15 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "backend.go"
Cohesion: 0.36
Nodes (4): askDecrypt(), pickDirsToDecrypt(), recordTick, hint

### Community 111 - "ConfinedEntry"
Cohesion: 0.13
Nodes (25): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache, ConfinedEntry() (+17 more)

### Community 112 - "installPage"
Cohesion: 0.07
Nodes (25): actionGate, binaryRow, binSource, handButton, handCheck, iconTile, installAction, installPage (+17 more)

### Community 113 - ".judgeProcess"
Cohesion: 0.15
Nodes (12): reasonText(), codeVouched(), isTraced(), loaderEnv(), readMountDevs(), roFor(), TestCodeVouched(), TestLoaderEnv() (+4 more)

### Community 115 - "maintainerKeeper"
Cohesion: 0.42
Nodes (4): appRunning(), libDirSig(), maintainedApp, maintainerKeeper

### Community 116 - "downloadAndVerify"
Cohesion: 0.15
Nodes (18): installHooks(), staged(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), FetchLatest(), fetchPair(), fetchRelease() (+10 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "UI"
Cohesion: 0.12
Nodes (18): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), InstallMode, InstallRequest, Options (+10 more)

### Community 121 - "go_pkg_path_filepath"
Cohesion: 0.07
Nodes (32): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), auditEntry(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree() (+24 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "newTestTarget"
Cohesion: 0.21
Nodes (18): newExeResolver(), FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps() (+10 more)

### Community 130 - "catalogRefresher"
Cohesion: 0.12
Nodes (13): drain(), readEvents(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad() (+5 more)

### Community 131 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Finder"
Cohesion: 0.14
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "serve.go"
Cohesion: 0.06
Nodes (36): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+28 more)

### Community 144 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 145 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.05
Nodes (68): newChangelogModel(), class, lineStyles, pos, viewportState, diffModel, editorModel, classOf() (+60 more)

### Community 146 - "Ledger"
Cohesion: 0.05
Nodes (39): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), candidates(), confirmAll(), describe(), record(), run() (+31 more)

### Community 147 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 148 - "GuardEvent"
Cohesion: 0.04
Nodes (24): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary(), GuardEvent, logBacklogDenial() (+16 more)

### Community 149 - "buildTrustedSet"
Cohesion: 0.10
Nodes (20): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+12 more)

### Community 150 - "RunWith"
Cohesion: 0.18
Nodes (19): RunWith(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError() (+11 more)

### Community 152 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 154 - "buildOneGuard"
Cohesion: 0.08
Nodes (34): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), displayEntries() (+26 more)

### Community 155 - "watchSet"
Cohesion: 0.09
Nodes (24): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), fakeInotify, watchDir (+16 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "sources"
Cohesion: 0.24
Nodes (5): TestGeneratedLibraryBlockRoundTrips(), pathByte(), pathRecord(), TestPathRecord(), sources

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "bufio.Reader"
Cohesion: 0.18
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 162 - "NetEventType"
Cohesion: 0.08
Nodes (15): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+7 more)

### Community 163 - "runOfflineEdit"
Cohesion: 0.16
Nodes (12): auditAfterEdit(), runOfflineEdit(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), TestDecryptDirectoriesEmpty(), VerifyEncryptedKeys(), ScanEncryptedCatalogDirs() (+4 more)

### Community 164 - "Active"
Cohesion: 0.21
Nodes (10): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), EditText(), Active(), Field (+2 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (60): fakeLink, pageDef, pageID, reel, resizeWatch, shell, snapEngines, snapNet (+52 more)

### Community 167 - "sources_test.go"
Cohesion: 0.20
Nodes (9): SystemTrusted(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers() (+1 more)

### Community 168 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 169 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 170 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "dialLiveSession"
Cohesion: 0.11
Nodes (20): TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+12 more)

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (32): lsmHook, statFunc, trustHook, trustRows, ensureBtrfsLayout(), deleteResKeys(), multicallFamily(), singleBinaryPaths() (+24 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "selection_test.go"
Cohesion: 0.26
Nodes (8): clipboardText(), setClipboard(), cellOf(), repeat(), TestSelectionEditing(), TestSelectionMouseDrag(), TestSelectionReadOnly(), TestSelectionShiftArrowsCopy()

### Community 178 - "deploy"
Cohesion: 0.12
Nodes (26): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI(), fetchGUI(), guiAssetsFor(), hasAsset() (+18 more)

### Community 179 - "IntegrationSuite"
Cohesion: 0.24
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 180 - "time.Time"
Cohesion: 0.27
Nodes (5): LockedError, sweepID, sweepMemo, Progress, Seen

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (35): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+27 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 195 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "recoverFileInPlace"
Cohesion: 0.19
Nodes (19): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+11 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate_test.go"
Cohesion: 0.11
Nodes (30): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), pinOwnerLikelyAlive(), readPinState() (+22 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.21
Nodes (17): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+9 more)

### Community 203 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 204 - "Recorder"
Cohesion: 0.15
Nodes (16): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+8 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (30): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+22 more)

### Community 209 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 210 - "bpfstats/main.go"
Cohesion: 0.40
Nodes (9): chownToSudoUser(), compare(), dumpFullLog(), dumpTrace(), loadOne(), main(), measure(), resizeMaps() (+1 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "KernelDev"
Cohesion: 0.12
Nodes (22): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), binaryStatOf(), KernelDev(), StatBinary() (+14 more)

### Community 213 - "migrate_test.go"
Cohesion: 0.25
Nodes (8): stripEntryFlags(), stripImmutableFlags(), TestCopyXattrsPreservesUserAttributes(), TestStripImmutableFlagsDoesNotHangOnFIFO(), TestStripImmutableFlagsSkipsDanglingSymlink(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 218 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 219 - ".run"
Cohesion: 0.21
Nodes (6): peerCred(), TestClientAuth(), TestSASL(), clientAuth(), serverAuth(), saslServer

### Community 220 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 221 - "Gate"
Cohesion: 0.15
Nodes (6): userGate, Caller, Gate, Role, Listen(), Judge

### Community 222 - "ValidateTargets"
Cohesion: 0.19
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), Existing (+1 more)

### Community 223 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 224 - "reloadOnce"
Cohesion: 0.09
Nodes (28): loadDaemonConfig(), reloadOnce(), reloadSlotsNeeded(), setConfinedHomes(), startGuardedDaemon(), startGuardedDaemonAbortable(), unlockPendingGroupRoots(), TestPendingGroupUnlockStopRetiresDrainOnce() (+20 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (49): recordTickCmd(), recordModel, listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts() (+41 more)

### Community 228 - "mountTable"
Cohesion: 0.15
Nodes (8): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), mountTable, editorHarness

### Community 229 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 231 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 232 - "procRec"
Cohesion: 0.33
Nodes (5): judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), procRec

### Community 233 - "Question"
Cohesion: 0.40
Nodes (3): Confirm(), Question, recordUI

### Community 235 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (64): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), libraryBlocksFromCandidates(), sectionsFromCandidates(), sharedSection(), discordCatalogEntry(), groupedDiscordConf() (+56 more)

### Community 239 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 240 - "shQuote"
Cohesion: 0.17
Nodes (5): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), IntegrationSuite

### Community 242 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 243 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 244 - ".showSized"
Cohesion: 0.09
Nodes (17): configView, installUI, sizedDialog, scrollable(), selectable(), fadeScrollOf(), newFadeScroll(), ask() (+9 more)

### Community 249 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "plan_test.go"
Cohesion: 0.27
Nodes (14): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+6 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 663 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **42 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `parseDaemonEvents`, `go_pkg_strings`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.015921931176168466 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `reloadOnce`, `daemon/daemon.go`, `GuardInodeKey`, `time.Duration`, `engine`, `BinaryEntry`, `DaemonEvent`, `guardUnitTest`, `.resKey`, `os.File`, `time.Time`, `GuardEvent`, `Resource`, `loginServer`, `daemonconfig.go`, `buildOneGuard`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09263157894736843 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `How it works`, `Installation`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06148088746324512 - nodes in this community are weakly interconnected._