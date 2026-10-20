# Graph Report - app-listener  (2026-10-20)

## Corpus Check
- 511 files · ~1,726,933 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6653 nodes · 29388 edges · 239 communities (195 shown, 44 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2717 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `adadb8fe`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- wipe
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- go_pkg_path_filepath
- FileTarget
- protectedPage
- github.com/charmbracelet/bubbles/textarea.Model
- render.go
- GuardInodeKey
- absPath
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- resolveBits
- guard.bpf.c
- TrustGuard
- Finder
- highlight_test.go
- prepareDaemonStart
- Candidate
- TestMonitorUseCaseLifecycle
- UI
- runNetworkMonitor
- picker
- root.go
- IntegrationSuite
- Message
- controlServer
- TempBinary
- NetEventType
- IntegrationSuite
- IntegrationSuite
- monitorUnitTest
- go_pkg_strings
- fileedit_model_test.go
- CandidateDir
- os.File
- fadeScroll
- systemd_test.go
- fyne.io/fyne/v2/driver/desktop.Cursor
- Client
- Config
- fyne.io/fyne/v2/widget.Label
- Plan
- planUpdaters
- mutedText
- engine
- go_pkg_os
- New
- Recorder
- exe_supersede.h
- updatePage
- RefreshSection
- User
- wm_test.go
- IntegrationSuite
- audit.go
- update_test.go
- BinaryEntry
- journal
- Guard
- fcntl
- bpfstats/main.go
- policyEditor
- codecheck/main.go
- rcUser
- watchSet
- process_vm_readv.c
- GuardEvent
- LibraryClosure
- SelectOrphans
- githubRelease
- net_tester/main.go
- run
- recordFlow
- .showSized
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- EventType
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- pinstate_test.go
- Well-known bypass classes
- hashCache
- jit_provenance.c
- inspector_probe.c
- recordops.c
- codeedit.go
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- batchInodeRows
- mountTable
- github.com/charmbracelet/bubbletea.Cmd
- NetGuardEvent
- IntegrationSuite
- github.com/spf13/cobra.Command
- find_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- newTestTarget
- UnifiedDiff
- trace-app-libs.sh
- IntegrationSuite
- vaultfs.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- runDaemon
- serve.go
- dialog_test.go
- appTheme
- Ledger
- Vault
- .step
- gui/session.go
- NetworkMonitorRepository
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- IntegrationSuite
- configedit_test.go
- catalogRefresher
- Installation
- backend_test.go
- graphify-refresh.sh
- .Layout
- reloadOnce
- .populateWatchInodes
- feed.go
- shell
- editView
- time.Time
- WithBottomBar
- newTestVetter
- sync.Mutex
- Limitations
- Troubleshooting
- IntegrationSuite
- runMonitor
- netinject.c
- CopyTreeWithProgress
- inode_dev.h
- Binary
- Compatibility
- How it works
- ValidateTargets
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- downloadFile
- NewEventFanout
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- update/update.go
- pathCache
- runOfflineEdit
- DaemonEvent
- changelog_test.go
- Daemon configuration
- uninstall/uninstall.go
- Desktop GUI
- IntegrationSuite
- .FilterExistingWhitelist
- options.go
- .eventTypeAllowed
- bleedLayout
- procstats.go
- PlanDaemonUnitUpdate
- .awaitDaemonUpOrDead
- TestCleanRel
- KernelDev
- tempGrantSetup
- Vault
- fakeVault
- ptrace_race.c
- prctl_gate
- fakeDaemon
- stripTempJournal
- supersede_probe.c
- IntegrationSuite
- parseGuardEvent
- loadDaemonConfig
- Classify
- github.com/Virgula0/app-listener/tools/codecheck
- harnessSuite
- editorModel
- taint_victim.c
- parsePasswd
- changelogModel
- .TestDaemon_Btrfs_CopyIoctlGate
- .TestDaemon_Bypass_PrctlExeFileSwap
- TestMain

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 113 edges
3. `GuardInodeKey` - 100 edges
4. `shell` - 91 edges
5. `Load()` - 90 edges
6. `absPath()` - 88 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (239 total, 44 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.18
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (134): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+126 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "wipe"
Cohesion: 0.13
Nodes (17): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (76): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+68 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (27): addErrKind, bpfGuardEvent, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+19 more)

### Community 9 - "mkdirs"
Cohesion: 0.09
Nodes (42): TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+34 more)

### Community 10 - "go_pkg_path_filepath"
Cohesion: 0.37
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 11 - "FileTarget"
Cohesion: 0.13
Nodes (10): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+2 more)

### Community 12 - "protectedPage"
Cohesion: 0.07
Nodes (15): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview, recordView (+7 more)

### Community 13 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.19
Nodes (9): pos, dedent(), lines(), repeatKey(), reposition(), cellToRune(), cursorPos(), Selection (+1 more)

### Community 14 - "render.go"
Cohesion: 0.15
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 15 - "GuardInodeKey"
Cohesion: 0.06
Nodes (31): keptEntries, refusedReplacements, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), newVerifyState(), supersedeUnnamed(), appletTag() (+23 more)

### Community 16 - "absPath"
Cohesion: 0.08
Nodes (8): IntegrationSuite, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, absPath()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.06
Nodes (13): daemonPage, logGrid, logHook, logLine, runHost, runView, journalLine(), TestJournalLine() (+5 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (44): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+36 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (13): ownerPaths(), findResource(), Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots() (+5 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "resolveBits"
Cohesion: 0.43
Nodes (4): statFunc, resolveBits(), statKey(), UpdaterPlan

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "TrustGuard"
Cohesion: 0.06
Nodes (18): lsmHook, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteInoKeys(), deleteResKeys() (+10 more)

### Community 27 - "Finder"
Cohesion: 0.15
Nodes (10): findButton, span, column(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 28 - "highlight_test.go"
Cohesion: 0.15
Nodes (19): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+11 more)

### Community 29 - "prepareDaemonStart"
Cohesion: 0.13
Nodes (17): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), confirmMasterKeyOverwrite(), newPinGeneration(), prepareDaemonStart(), runBPFCheck(), runGenKey() (+9 more)

### Community 30 - "Candidate"
Cohesion: 0.22
Nodes (14): selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates(), sectionsFromCandidates() (+6 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.13
Nodes (14): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+6 more)

### Community 32 - "UI"
Cohesion: 0.09
Nodes (19): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, ShowDiff() (+11 more)

### Community 33 - "runNetworkMonitor"
Cohesion: 0.09
Nodes (18): runNetworkMonitor(), runTUI(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), appendCapped(), clipToWidth() (+10 more)

### Community 34 - "picker"
Cohesion: 0.09
Nodes (22): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+14 more)

### Community 35 - "root.go"
Cohesion: 0.10
Nodes (18): installHooks(), releaseBuild(), Run(), staged(), closeDumpFile(), Execute(), initLogger(), postRun() (+10 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (12): checkSecret(), TestLiveSessionWatchAndPinger(), File, backend, backendOp, backendState, fakeLink, BinaryCheck (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (24): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+16 more)

### Community 39 - "TempBinary"
Cohesion: 0.07
Nodes (25): confirmUserPlaced(), forwardCheck(), forwardRequest, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+17 more)

### Community 40 - "NetEventType"
Cohesion: 0.06
Nodes (24): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), BinariesSummary(), TestBinariesSummary() (+16 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "IntegrationSuite"
Cohesion: 0.10
Nodes (4): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite

### Community 44 - "go_pkg_strings"
Cohesion: 0.06
Nodes (14): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin() (+6 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "CandidateDir"
Cohesion: 0.08
Nodes (27): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+19 more)

### Community 48 - "os.File"
Cohesion: 0.04
Nodes (75): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir() (+67 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "systemd_test.go"
Cohesion: 0.05
Nodes (68): offerDaemonStart(), runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+60 more)

### Community 51 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.11
Nodes (12): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), newButtonCursor(), newCheckGroup(), newRadioGroup() (+4 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (15): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+7 more)

### Community 53 - "Config"
Cohesion: 0.07
Nodes (56): swapFile(), addedResources(), appendSectionsAndEdit(), collectDiffAdditions(), TestAddedResourcesSkipsAlreadyGuardedSSH(), TestUncoveredCandidates(), uncoveredCandidates(), grantNewBinaries() (+48 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (24): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+16 more)

### Community 55 - "Plan"
Cohesion: 0.12
Nodes (15): Plan, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded(), TestPlanExtendsAnExistingSection() (+7 more)

### Community 56 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.05
Nodes (45): actionGate, binSource, checkGrid, filesPage, handCheck, homePage, installAction, installPage (+37 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (28): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), pathWithin() (+20 more)

### Community 59 - "go_pkg_os"
Cohesion: 0.10
Nodes (34): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+26 more)

### Community 60 - "New"
Cohesion: 0.08
Nodes (50): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+42 more)

### Community 61 - "Recorder"
Cohesion: 0.21
Nodes (12): Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), TestMonitorTargetsNameFilesByAbsolutePath(), TestProgressNamesMulticallApplets() (+4 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "updatePage"
Cohesion: 0.12
Nodes (13): anyDaemonRunning(), installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), newUpdatePage() (+5 more)

### Community 64 - "RefreshSection"
Cohesion: 0.10
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 65 - "User"
Cohesion: 0.06
Nodes (71): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), bunLaunchersOffered() (+63 more)

### Community 66 - "wm_test.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "audit.go"
Cohesion: 0.25
Nodes (14): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+6 more)

### Community 69 - "update_test.go"
Cohesion: 0.11
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 70 - "BinaryEntry"
Cohesion: 0.09
Nodes (18): BinaryEntry, diffBinaries(), TestDiffBinaries(), netGuardPill(), entryOf(), EntryWithHash(), NetworkGuardRepository, newExeResolver() (+10 more)

### Community 71 - "journal"
Cohesion: 0.14
Nodes (8): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 72 - "Guard"
Cohesion: 0.06
Nodes (18): deferredBinary, exeRow, rootHandle, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), confinedEntry() (+10 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "policyEditor"
Cohesion: 0.13
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+5 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "rcUser"
Cohesion: 0.12
Nodes (41): TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent() (+33 more)

### Community 78 - "watchSet"
Cohesion: 0.24
Nodes (6): isSymlink(), watchDir, watchPos, watchSet, TestNearestRootClaims(), nearestRootClaims()

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.05
Nodes (11): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary(), GuardEvent (+3 more)

### Community 81 - "LibraryClosure"
Cohesion: 0.07
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 82 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 83 - "githubRelease"
Cohesion: 0.27
Nodes (16): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+8 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "run"
Cohesion: 0.24
Nodes (11): deployGUI(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), RunWith(), syncGUI() (+3 more)

### Community 86 - "recordFlow"
Cohesion: 0.13
Nodes (20): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickBinaries(), pickEncryption() (+12 more)

### Community 87 - ".showSized"
Cohesion: 0.09
Nodes (15): binaryList, installUI, multiAnswer, sizedDialog, scrollable(), selectable(), showErrorDialog(), newFadeScroll() (+7 more)

### Community 88 - "liveSession"
Cohesion: 0.09
Nodes (29): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+21 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.06
Nodes (22): systemPatterns(), systemRule(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents() (+14 more)

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (11): PumpRingbuf(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+3 more)

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
Nodes (13): procStats, statsSample, formatClock(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+5 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.12
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 101 - "pinstate_test.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "hashCache"
Cohesion: 0.12
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "codeedit.go"
Cohesion: 0.16
Nodes (24): class, diffModel, classOf(), Edit(), leadingIndent(), moveTo(), Navigate(), New() (+16 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 111 - "mountTable"
Cohesion: 0.14
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), mountTable (+1 more)

### Community 112 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (34): recordTickCmd(), recordModel, listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), listenForGuardEvents(), tickGuardStats(), listenForEvents() (+26 more)

### Community 113 - "NetGuardEvent"
Cohesion: 0.11
Nodes (15): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), snapNetGuard (+7 more)

### Community 115 - "github.com/spf13/cobra.Command"
Cohesion: 0.06
Nodes (48): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+40 more)

### Community 116 - "find_test.go"
Cohesion: 0.21
Nodes (18): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+10 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (64): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), confirmRunPrereq(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+56 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 131 - "UnifiedDiff"
Cohesion: 0.40
Nodes (6): refreshCatalogConfig(), installConfig(), ConfirmOverwrite(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "vaultfs.go"
Cohesion: 0.05
Nodes (44): entryOf(), readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile() (+36 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "runDaemon"
Cohesion: 0.10
Nodes (26): makeReloadHandler(), noteReloadedSettings(), runDaemon(), startCatalogRefresh(), startPprof(), writePidFile(), inspectorPaths(), resolveInspectors() (+18 more)

### Community 143 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 144 - "dialog_test.go"
Cohesion: 0.41
Nodes (11): openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder(), TestPickerSubfolderSearchPicksFile() (+3 more)

### Community 146 - "Ledger"
Cohesion: 0.08
Nodes (19): applySchema(), EnsurePlaceholders(), Ledger, Pending, Source, JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir() (+11 more)

### Community 147 - "Vault"
Cohesion: 0.18
Nodes (16): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), resolveFilesystemPrereqs(), askEncryption() (+8 more)

### Community 148 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 149 - "gui/session.go"
Cohesion: 0.21
Nodes (14): ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestOwnedByChecksTypeAndOwner(), TestReadOwnedFileRefusesSymlink() (+6 more)

### Community 150 - "NetworkMonitorRepository"
Cohesion: 0.15
Nodes (6): runHeadless(), ProtocolString(), NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (43): DrainUntilSignal(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards() (+35 more)

### Community 156 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 157 - "catalogRefresher"
Cohesion: 0.10
Nodes (20): drain(), newCatalogRefresher(), readEvents(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan() (+12 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 163 - "reloadOnce"
Cohesion: 0.22
Nodes (13): backingDeviceUnion(), buildGuards(), drainEphemeral(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), startGuardedDaemonAbortable(), unlockPendingGroupRoots() (+5 more)

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (38): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+30 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (44): pageDef, pageID, reel, resizeWatch, shell, snapEngines, snapNet, sysStatus (+36 more)

### Community 167 - "editView"
Cohesion: 0.08
Nodes (15): codeEntry, editView, Entry, handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm(), showDialog() (+7 more)

### Community 168 - "time.Time"
Cohesion: 0.21
Nodes (7): LockedError, newWatchSet(), fakeInotify, sweepID, sweepMemo, Progress, Seen

### Community 169 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 171 - "sync.Mutex"
Cohesion: 0.22
Nodes (6): newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "runMonitor"
Cohesion: 0.12
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), pathWithinDir() (+11 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "Binary"
Cohesion: 0.14
Nodes (17): lines(), pickLibs(), beyondEvents(), finishBinary(), Binary, isShellHelper(), newBinary(), noteWhitelisted() (+9 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "ValidateTargets"
Cohesion: 0.20
Nodes (8): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), Existing, ExistingRule

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 195 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "update/update.go"
Cohesion: 0.11
Nodes (28): applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), downloadAndVerify() (+20 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "runOfflineEdit"
Cohesion: 0.13
Nodes (15): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+7 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.12
Nodes (31): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+23 more)

### Community 203 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 204 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 205 - "uninstall/uninstall.go"
Cohesion: 0.09
Nodes (39): deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), confirmKeyDeletion(), decryptDirectories(), decryptStep(), offerBackupCleanup() (+31 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - ".FilterExistingWhitelist"
Cohesion: 0.09
Nodes (23): confSafeMatch(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed() (+15 more)

### Community 210 - "options.go"
Cohesion: 0.36
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 214 - "PlanDaemonUnitUpdate"
Cohesion: 0.27
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 217 - "KernelDev"
Cohesion: 0.09
Nodes (25): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+17 more)

### Community 218 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 220 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 224 - "stripTempJournal"
Cohesion: 0.32
Nodes (8): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal()

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 227 - "parseGuardEvent"
Cohesion: 0.29
Nodes (8): decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 228 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 229 - "Classify"
Cohesion: 0.40
Nodes (4): Classify(), drain(), TestClassifyVerifiesTheInodeThatRan(), harnessSuite

### Community 234 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 240 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 616 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **44 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `uninstall/uninstall.go`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.020289569585344232 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestDaemon_Btrfs_CopyIoctlGate`, `.TestDaemon_Bypass_PrctlExeFileSwap`, `IntegrationSuite`, `go_pkg_strings`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `wipe` be split into smaller, more focused modules?**
  _Cohesion score 0.13118279569892474 - nodes in this community are weakly interconnected._